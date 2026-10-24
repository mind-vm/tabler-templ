package auth_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/avatar"
	"github.com/mind-vm/tabler-templ/blocks/auth"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestSignInCard(t *testing.T) {
	h := render(t, auth.SignInCard(auth.SignInCardProps{
		Form: auth.SignInFormProps{Action: "/login", ForgotHref: "/forgot", RememberName: "remember", Error: "Wrong password", UserValue: "ann@example.com",
			Hidden: []auth.Hidden{{Name: "next", Value: "/orders"}}},
		Providers: []auth.Provider{{Href: "/auth/github", Icon: "ti-brand-github", Label: "Login with Github"}},
	}))
	assert.Contains(t, h, `<form action="/login" method="post" novalidate><input type="hidden" name="next" value="/orders">`)
	assert.Contains(t, h, `Wrong password`)
	assert.Contains(t, h, `value="ann@example.com" autocomplete="email" required aria-invalid="true" id="signin-email" name="email"`)
	assert.Contains(t, h, `<a href="/forgot" class="form-label-description ms-auto">I forgot password</a>`)
	assert.Contains(t, h, `data-password-toggle="signin-password"`)
	assert.Contains(t, h, `<div class="hr-text">or</div>`)
	assert.Contains(t, h, `<a href="/auth/github" class="btn w-100">`)
}

func TestCards(t *testing.T) {
	h := render(t, auth.SignUpCard(auth.SignUpCardProps{Action: "/signup", TermsHref: "/terms"}))
	assert.Contains(t, h, `data-bs-input="#signup-password"`)
	assert.Contains(t, h, `<a href="/terms">terms and policy</a>.`)
	// The toggle's script must not sit between the field and the meter (Tabler's sibling spacing).
	assert.Less(t, strings.Index(h, `class="strength"`), strings.Index(h, `<script>`))

	h = render(t, auth.VerifyCard(auth.VerifyCardProps{Action: "/verify"}))
	assert.Contains(t, h, `<div class="otp otp-lg" data-bs-toggle="otp" data-bs-length="6" data-bs-groups="[3,3]"><input type="text" name="code" maxlength="6" aria-label="Verification code"></div>`)
	assert.Contains(t, render(t, auth.LockCard(auth.LockCardProps{Name: "Ann", Avatar: avatarProps()})), `class="avatar avatar-xl mb-3"`)
}

func TestPages(t *testing.T) {
	h := render(t, auth.Page(auth.PageProps{Title: "Sign in", Lang: "de"}))
	assert.Contains(t, h, `<html lang="de">`)
	assert.Contains(t, h, `<main class="page page-center" id="content"><div class="container container-tight py-4">`)
	h = render(t, auth.CoverPage(auth.PageProps{}, auth.CoverProps{Image: "/cover.jpg"}))
	assert.Contains(t, h, `<body class="d-flex flex-column bg-surface">`)
	assert.Contains(t, h, `style="background-image: url(/cover.jpg);"`)
}

func avatarProps() avatar.Props { return avatar.Props{Initials: "AN"} }
