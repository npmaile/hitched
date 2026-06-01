# Wedding Reception Website — Edit Guide

This is a Go-powered website for Nate & Susanna Maile's wedding reception. The Go binary embeds all HTML at compile time, so you deploy a single executable with an SQLite database file next to it.

---

## Quick facts

| Detail | Value |
|---|---|
| Names | Nate & Susanna Maile |
| Date | August 15, 2025 |
| Time | 3:00 PM |
| Venue | The Supermarket |
| Address | 638 N. Highland Ave NE, Atlanta, GA 30306 |
| Coordinates | 33.77163534989877, -84.35295692457922 |
| Dress code | Casual |
| RSVP deadline | August 1st |

---

## Project structure

```
hitched/
├── main.go          # Go HTTP server
├── go.mod / go.sum  # dependencies (modernc.org/sqlite)
├── index.html       # Home / RSVP page
├── hotels.html      # Hotels & lodging page
├── registry.html    # Art registry page
└── rsvps.db         # SQLite database (created on first run, not in repo)
```

---

## Running the server

```bash
# Build
go build -o hitched .

# Run with defaults (:8080, ./rsvps.db)
./hitched

# Custom port and database path
./hitched -addr :9000 -db /var/data/rsvps.db
```

The server sits behind an nginx reverse proxy for TLS — no TLS configuration is needed in the binary itself.

If the process crashes and leaves a corrupt `rsvps.db` (zero-byte file), delete it and restart. SQLite will create a fresh one.

### Querying RSVPs

```bash
sqlite3 rsvps.db "SELECT first_name, last_name, email, attending, additional_guests, submitted_at FROM rsvps ORDER BY created_at;"
```

---

## Design system

Bold, graphic aesthetic with chunky 2.5px black borders and flat color blocks. No shadows, no gradients.

### Color palette (CSS variables in `:root` of each file)

| Variable | Hex | Used for |
|---|---|---|
| `--lime` | `#c8f135` | Marquee banner, form focus state, nav active link |
| `--coral` | `#ff5e3a` | Hero right panel, eyebrow text, submit button, registry header |
| `--sky` | `#4fc3f7` | Details card 1, hotels page header |
| `--violet` | `#b388ff` | RSVP section background |
| `--pink` | `#ff80ab` | Details card 3, decorative shapes |
| `--yellow` | `#ffe033` | Hover states, venue map panel, row hover on registry |
| `--white` | `#fffef9` | Page background, form boxes |
| `--ink` | `#1a1008` | All borders, text, nav/footer backgrounds |

All four files share the same `:root` block. To change a color, update it in each file (or consider extracting to a shared CSS file if you add more pages).

### Fonts (loaded from Google Fonts)

- **Fraunces** — all display headings, card titles, prices, hero names
- **Space Grotesk** — body text, labels, buttons, form fields

---

## Pages

### index.html — Home / RSVP

Sections top to bottom:

1. **Sticky nav bar** — links to all three pages. `active` class on current page.
2. **Marquee banner** — scrolling ticker. Text is duplicated inside `.marquee-inner` for the seamless loop; keep both halves identical when editing.
3. **Hero** — 50/50 split: names on the left, big date on coral background on the right. Right panel hidden on mobile.
4. **Details cards** — three colored cards (date, venue, dress code) in a CSS grid.
5. **Venue map** — two-column section: address/links on the left (`--yellow` background), OpenStreetMap iframe on the right. Map coordinates are set in the iframe `src` — see [Updating the map pin](#updating-the-map-pin).
6. **RSVP form** — see [RSVP form](#rsvp-form) below.
7. **Footer** — dark background with lime name display.

### hotels.html — Hotels & Lodging

- Page header with `--sky` background
- Full-width OpenStreetMap iframe showing the area around the venue (wider zoom than the venue map)
- Grid of hotel cards, each with name, address, distance, and a Google Maps link
- Note band at the bottom with booking advice

### registry.html — Art Registry

- Split header: message on coral background, big palette emoji on pink background
- List of 14 art ideas ordered from easiest (~$0) to most ambitious (~$250)
- Each row has a title, description, price, and a color-coded difficulty badge (lime = easy, yellow = medium, coral = hard)
- Row hover highlights in `--yellow`

---

## RSVP form

The form collects: first name, last name, email, additional guests, attendance choice, and notes.

### Attendance options

There are 9 options in a 3-column grid, each with a `data-val` of `yes`, `no`, or `maybe`:

| Label | Value |
|---|---|
| Joyfully accept | yes |
| Ecstatically accept | yes |
| Reluctantly accept | yes |
| Regretfully accept | yes |
| Desperately accept | yes |
| There in spirit | maybe |
| Regretfully decline | no |
| Tearfully decline | no |
| Gleefully decline | no |

### Emoji animation

When the submit button is clicked, `launchEmojis(attendance)` fires before the success state appears:

- `yes` — confetti burst (🎉🎊✨🌟💛🎈🥂💐) radiating outward from random screen positions
- `no` — emoji rainstorm (🌧️💧☁️😢🥀⛈️😭💔) falling from the top
- `maybe` — ghosts and wisps (👻✨💫🌫️🕯️🌙) rising from the bottom

Animations clean up automatically after 5 seconds.

### Success messages

The success message and emoji vary by which button was selected. Funny choices (gleefully decline, reluctantly accept, etc.) get their own copy. See `showSuccess()` in the `<script>` block of `index.html`.

### Backend

The form POSTs JSON to `POST /api/rsvp`. Payload shape:

```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane@example.com",
  "additional_guests": 1,
  "attending": "yes",
  "dietary_notes": "Vegetarian",
  "submitted_at": "2025-07-15T14:30:00.000Z"
}
```

The server validates that `first_name` and `email` are present and returns `201` on success. The frontend `catch` block currently shows the success state even on network failure — remove that fallback once you've confirmed the backend is reachable from production.

---

## Updating the map pin

Both maps use OpenStreetMap embed iframes. The pin is set via the `marker` query parameter. Current coordinates: `33.77163534989877, -84.35295692457922`.

**Venue map** (index.html):
```
bbox=-84.3630,33.7666,-84.3430,33.7766&layer=mapnik&marker=33.77163...,−84.35295...
```
The `bbox` defines the visible area — it's centered on the pin. If you move the pin significantly, update both `bbox` and `marker`.

**Hotels map** (hotels.html):
```
bbox=-84.4050,33.7400,-84.3250,33.8050&layer=mapnik&marker=33.77163...,-84.35295...
```
The bbox here is intentionally wide to show the surrounding neighborhood. Only `marker` needs to change if the venue moves.

---

## Common edits

**Change the date/time:**
Search all HTML files for `August 15` and `3:00 PM` and update. The marquee text in `index.html` is duplicated — update both halves.

**Change the venue or address:**
Update the text in the details card, the venue map panel, the footer in all three files, the Google Maps/Apple Maps `href` values in the map section, and the `marker`/`bbox` values in both iframe `src` attributes.

**Add a hotel:**
Copy an existing `.hotel-card` div in `hotels.html` and update the name, address, distance, and Google Maps link. The `:nth-child()` color pattern cycles through the palette — adding a 7th card will re-use `--lime`.

**Add an art idea to the registry:**
Copy an existing `.art-item` block in `registry.html`, update the number, title, description, emoji, price, and difficulty class (`diff-easy`, `diff-medium`, or `diff-hard`).

**Change a color globally:**
Update the hex value in `:root` in the relevant file. All components in that file inherit from those variables.

**Remove the marquee:**
Delete the `<div class="marquee-wrap">` block and its CSS rules from `index.html`.
