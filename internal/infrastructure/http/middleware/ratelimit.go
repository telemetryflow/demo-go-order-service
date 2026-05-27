// Package middleware provides HTTP rate limiting middleware.
//
// TelemetryFlow Order Service - Community Enterprise Observability Platform
// Copyright (c) 2024-2026 Telemetri Data Indonesia. All rights reserved.
// Open Source Software built by Telemetri Data Indonesia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/telemetryflow/order-service/internal/infrastructure/config"
)

// rateLimiter implements a simple token bucket rate limiter
type rateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientBucket
	requests int
	window   time.Duration
}

type clientBucket struct {
	tokens    int
	lastReset time.Time
}

// newRateLimiter creates a new rate limiter
func newRateLimiter(requests int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		clients:  make(map[string]*clientBucket),
		requests: requests,
		window:   window,
	}

	// Cleanup old entries periodically
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// allow checks if a request is allowed
func (rl *rateLimiter) allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	client, exists := rl.clients[clientIP]
	if !exists {
		rl.clients[clientIP] = &clientBucket{
			tokens:    rl.requests - 1,
			lastReset: now,
		}
		return true
	}

	// Reset tokens if window has passed
	if now.Sub(client.lastReset) > rl.window {
		client.tokens = rl.requests - 1
		client.lastReset = now
		return true
	}

	// Check if tokens available
	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false
}

// cleanup removes old entries
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, client := range rl.clients {
		if now.Sub(client.lastReset) > rl.window*2 {
			delete(rl.clients, ip)
		}
	}
}

// RateLimit returns rate limiting middleware
func RateLimit(cfg config.RateLimitConfig) echo.MiddlewareFunc {
	limiter := newRateLimiter(cfg.Requests, cfg.Window)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientIP := c.RealIP()

			if !limiter.allow(clientIP) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}

			return next(c)
		}
	}
}
