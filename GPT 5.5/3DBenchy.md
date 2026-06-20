# Animated 3D Benchy SVG Prompt

## English Version

### Prompt

Create a fully self-contained SVG file featuring an animated pseudo-3D **3D Benchy boat** — the iconic benchmark model used for testing 3D printing quality.

The boat must remain visually recognizable as the classic 3D Benchy benchmark model, with accurate tugboat proportions and benchmark-test details.

The SVG should present the boat in an isometric or slightly perspective 3D view, with clean vector geometry, no raster images, no external assets, no JavaScript, and only inline SVG, CSS, and/or SMIL animations.

---

### Technical Requirements

* Return only the complete SVG code, starting with `<svg ...>` and ending with `</svg>`.
* Use `viewBox="0 0 1200 800"`.
* The SVG must be scalable, cleanly structured, and browser-compatible.
* Group all major parts using clear `<g>` IDs:

  * `benchy`
  * `hull`
  * `cabin`
  * `chimney`
  * `windows`
  * `deck`
  * `waves`
  * `shadow`
  * `highlights`
  * `details`
* Include `<title>` and `<desc>` for accessibility.
* Use an inline `<style>` block for CSS animations.
* Do not use external fonts, images, links, scripts, or dependencies.
* Use vector shapes only:

  * `<path>`
  * `<polygon>`
  * `<ellipse>`
  * `<circle>`
  * `<rect>`
  * `<linearGradient>`
  * `<radialGradient>`
  * SVG filters such as `<filter>`, `<feGaussianBlur>`, and `<feDropShadow>`

---

### Visual Style

The visual style should feel like a polished technical 3D illustration:

* Crisp outlines
* Soft shadows
* Plastic-like gradients
* Bright edge highlights
* Clear pseudo-3D depth
* Modern vector illustration quality
* Detailed but not photorealistic

The result should be suitable for:

* A landing page
* A 3D printing tool
* A technical documentation page
* An animated SVG demo
* A product UI related to 3D printing

---

### Composition

Place the 3D Benchy boat in the center of the canvas.

Show the boat from a 3/4 perspective so the following parts are visible:

* Bow
* Starboard side
* Deck
* Cabin
* Chimney
* Part of the stern

The boat should immediately read as a 3D Benchy: a compact tugboat-like benchmark model with a rounded hull, central cabin, chimney, windows, portholes, and sturdy proportions.

---

### Hull Details

The hull should be short, solid, and compact.

Include:

* Rounded bow
* Slightly raised deck
* Thick side wall
* Smooth lower body
* Flat or slightly flattened stern
* Slightly protruding upper rim
* Smooth widening from the bow toward the center
* Visible lower curved hull section
* Inner shadow under the upper rim

The bow should have the characteristic Benchy shape:

* Rounded front
* Slightly protruding upper edge
* Smooth transition into the side wall
* Compact tugboat proportions

The stern should be flatter and slightly higher than the bottom hull line.

---

### Cabin and Chimney

Add a central cabin structure with:

* Slightly slanted walls
* A visible roof
* A large front window
* Side windows
* Thick roof edge
* Subtle 3D side planes

Add a chimney on top of the cabin.

The chimney should appear cylindrical in pseudo-3D and include:

* Elliptical top
* Dark inner opening
* Shaded side wall
* Gradient plastic surface
* Slight highlight on the upper rim

The chimney should move naturally with the boat animation.

Do not add separate smoke, so the illustration remains clean and technical.

---

### Windows and Portholes

Add 2–3 circular portholes on the side of the hull.

Each porthole should include:

* Dark inner fill
* Slight rim highlight
* Inner shadow
* Small glossy reflection

Cabin windows should be:

* Dark blue or almost black
* Slightly reflective
* Made from rounded rectangles or polygons
* Accented with subtle cyan-blue highlights

---

### 3D Printing Benchmark Details

Add small technical details associated with 3D printing test models:

* Subtle horizontal layer lines
* Overhang-like shapes
* Small arches
* Bridge-like elements
* Holes
* Sharp edges
* Rounded edges
* Tolerance-test-like micro details
* Thin contour lines
* Small engraved or embossed marks

Optionally add a small vector-style label on the hull, such as:

* `3DBenchy`
* `BENCHY TEST`

Do not rely on external fonts.

---

### Color and Material

The boat should look like glossy PLA plastic.

Use:

* Main color: rich blue, cyan-blue, or turquoise-blue
* Shadows: deep navy blue and dark teal
* Highlights: pale cyan, icy blue, and soft white-blue
* Window interiors: very dark blue, almost black
* Hole interiors: very dark blue, almost black

Use gradients on:

* Hull
* Cabin
* Roof
* Chimney
* Deck surfaces
* Side planes

Add thin bright highlights along:

* Upper hull rim
* Cabin roof
* Bow edge
* Chimney top
* Porthole rims
* Window reflections

Add a soft cast shadow under the boat as a blurred ellipse.

---

### Animation

The boat should gently rock on waves.

Animation requirements:

* Smooth side-to-side rotation of approximately 2–4 degrees
* Subtle vertical bobbing of 6–10 pixels
* Infinite animation loop
* Duration: approximately 3–4 seconds
* Timing: `ease-in-out`
* Motion should feel calm, smooth, and unobtrusive

The waves under the boat should move horizontally in a continuous loop.

Add 2–3 wave layers:

* Brighter foreground wave layer
* Softer middle wave layer
* More transparent background wave layer

Add a subtle animated glossy highlight on the hull, such as a narrow light streak slowly moving across the side surface.

Respect `prefers-reduced-motion`:

* Disable animations when reduced motion is enabled
* Or make the movement almost imperceptible

---

### Pseudo-3D Construction

Build the pseudo-3D boat from multiple visible surfaces.

The hull should include:

* Top deck plane
* Side surface
* Bow-facing surface
* Lower curved body
* Inner shadow below the rim

The cabin should include:

* Front wall
* Side wall
* Roof top
* Roof thickness
* Window planes

The chimney should include:

* Top ellipse
* Bottom ellipse
* Curved side wall
* Dark inner opening
* Gradient shading

Use SVG filters for depth:

* Drop shadow
* Gaussian blur
* Soft cast shadow
* Subtle inner-like shadows using layered paths

Outlines should be clean and consistent:

* Around 2–4 px
* Slightly darker than the fill color
* Smooth and polished

---

### Background

The background should be transparent or a very light technical gradient.

Prefer a transparent background, with water and waves forming a soft base beneath the boat.

Do not overload the background with scenery.

The scene should look good on both light and dark page backgrounds.

---

### Restrictions

Do not make the boat look like:

* A sailboat
* A yacht
* A military ship
* A container ship
* A generic cartoon boat

Avoid:

* Realistic rendering
* Pixel art
* Low-poly style
* Cartoon eyes
* Characters
* Humans
* Third-party logos
* Brand marks
* External images
* External fonts
* JavaScript

Keep the design technical, polished, playful, and precise.

---

### Final Output

Produce one complete self-contained SVG file.

The SVG code must be neatly formatted.

It must open directly in a modern browser and immediately display the animated pseudo-3D Benchy.

All animations must work without JavaScript or external libraries.

Before returning the result, verify that:

* The SVG is valid
* All tags are closed
* All IDs are unique
* All animations work without JavaScript
* The boat remains recognizable as the classic 3D Benchy benchmark model
