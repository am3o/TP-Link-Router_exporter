package router

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const session = "Session"

type Client struct {
	url   string
	token string
}

func New(url, user, password string) (Client, error) {
	hash := md5.New()
	if _, err := hash.Write([]byte(password)); err != nil {
		return Client{}, fmt.Errorf("could not create valid password: %w", err)
	}

	return Client{
		url: url,
		token: base64.StdEncoding.EncodeToString(
			fmt.Appendf([]byte{}, "%v:%v", user, hex.EncodeToString(hash.Sum(nil))),
		),
	}, nil
}

func (c *Client) Login(ctx context.Context) (context.Context, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%v/userRpm/LoginRpm.htm?Save=Save", c.url), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create login: %w", err)
	}

	req.Header.Set("Authorization", c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute login: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read login content: %w", err)
	}
	_ = content

	return context.WithValue(ctx, session, "TODO"), nil
}

func (c *Client) Logout(ctx context.Context) error {
	session := ctx.Value(session)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%v/%s/userRpm/LogoutRpm.htm", session, c.url), nil)
	if err != nil {
		return fmt.Errorf("could not logout: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not logout: %w", err)
	}
	_ = resp
	if resp.StatusCode != http.StatusOK {
		return errors.New("could not logout: router does not accept the logout")
	}

	return nil
}

func (c *Client) WANTraffic(ctx context.Context) (rx map[string]float64, tx map[string]float64, err error) {
	// TODO: parse WAN information from the router
	return
}

func (c *Client) LANTraffic(ctx context.Context) (rx map[string]float64, tx map[string]float64, err error) {
	// TODO: parse LAN information from the router
	return
}
