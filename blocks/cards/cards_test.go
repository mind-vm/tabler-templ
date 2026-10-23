package cards_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/avatar"
	"github.com/mind-vm/tabler-templ/blocks/cards"
	"github.com/mind-vm/tabler-templ/blocks/tabscard"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func with(c templ.Component, children string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(templ.WithChildren(ctx, templ.Raw(children)), w)
	})
}

func TestPlaceholder(t *testing.T) {
	h := render(t, cards.Placeholder(cards.PlaceholderProps{Height: 150}))
	assert.Contains(t, h, `<div class="card-body"><svg xmlns="http://www.w3.org/2000/svg" class="w-100 border rounded" preserveAspectRatio="none" width="400" height="150" viewBox="0 0 400 150"`)
	assert.Contains(t, h, `<line x1="0" y1="150" x2="400" y2="0"></line>`)
	assert.Contains(t, render(t, cards.Placeholder()), `height="200"`)
}

func TestBottomTabs(t *testing.T) {
	h := render(t, cards.BottomTabs(tabscard.Props{ID: "t"}, []tabscard.Tab{{Key: "a", Label: "A"}, {Key: "b", Label: "B"}}))
	assert.Contains(t, h, `<div id="t-a" class="card tab-pane active show" role="tabpanel" aria-labelledby="t-a-tab">`)
	assert.Contains(t, h, `nav-tabs-bottom`)
	assert.Less(t, strings.Index(h, "tab-content"), strings.Index(h, "nav-tabs-bottom"), "tabs below the panes")
}

func TestImages(t *testing.T) {
	assert.Equal(t, `<div class="img-responsive img-responsive-21x9 card-img-bottom" style="background-image: url(p.jpg);"></div>`,
		render(t, cards.Photo(cards.PhotoProps{Src: "p.jpg", Bottom: true})))
	h := render(t, with(cards.SideImage(cards.SideImageProps{Src: "p.jpg", Alt: "P", End: true}), "body"))
	assert.Contains(t, h, `<div class="row row-0"><div class="col-3 order-md-last"><img src="p.jpg" alt="P" class="card-img-end w-100 h-100 object-cover"></div><div class="col">body</div></div>`)
}

func TestUsersList(t *testing.T) {
	users := []cards.User{{Name: "Ann", Avatar: avatar.Props{Initials: "AN"}, Description: "Hi", Color: "red", Checked: true}, {Name: "Bob", Description: "Yo"}}
	h := render(t, cards.UsersList(cards.UsersListProps{Title: "Last commits", Hoverable: true}, users))
	assert.Contains(t, h, `list-group-hoverable`)
	assert.Contains(t, h, `<span class="badge bg-red"></span>`)
	assert.Contains(t, h, `<a href="#" class="list-group-item-title">Ann</a><div class="list-group-item-description">Hi</div>`)
	assert.Contains(t, h, `list-group-item-actions show`)
	assert.Contains(t, h, `aria-label="Remove from favorites"`)

	h = render(t, cards.UsersList(cards.UsersListProps{Title: "Contacts", Checkbox: true}, users))
	assert.Contains(t, h, `list-group-selectable`)
	assert.Contains(t, h, `<input type="checkbox" class="form-check-input list-group-selectable-check m-0" checked aria-label="Select Ann">`)
}

func TestPeopleAndTopUsers(t *testing.T) {
	h := render(t, cards.People(cards.PeopleProps{Title: "People", MaxHeight: "35rem"}, []cards.Group{{Letter: "A", Users: []cards.User{{Name: "Ann"}}}}))
	assert.Contains(t, h, `style="max-height: 35rem"`)
	assert.Contains(t, h, `<div class="list-group-header sticky-top">A</div>`)

	h = render(t, cards.TopUsers("Top users", []cards.User{{Name: "Ann", Status: "green", StatusLabel: "Online", Meta: "now"}}))
	assert.Contains(t, h, `<span class="badge bg-green" role="img" aria-label="Online"></span>`)
	assert.Contains(t, h, `<div class="text-secondary text-truncate mt-n1">now</div>`)
}

func TestInviteAndCreditCard(t *testing.T) {
	h := render(t, with(cards.Invite(cards.InviteProps{Title: "No Team Members", Avatars: []avatar.Props{{Initials: "A"}}}), "Invite"))
	assert.Contains(t, h, `card card-lg card-dashed card-transparent`)
	assert.Contains(t, h, `<p class="text-secondary">Invite</p>`)
	h = render(t, cards.CreditCard(cards.CreditCardProps{}))
	assert.Contains(t, h, `id="credit-card-number"`)
	assert.Contains(t, h, `<option value="2030">2030</option>`)
	assert.Contains(t, h, `Pay now`)
}
