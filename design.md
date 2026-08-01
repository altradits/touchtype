# GateCheck Design System (v1.0)

## 🎨 Design Philosophy: "Bare-Minimum Elegance"
The UI is designed for a **high-pressure, live event environment**. It prioritizes:
1. **Legibility:** Large fonts, high contrast, and generous whitespace so it can be read on a projector or a bright smartphone screen.
2. **Trustworthiness:** A clean, sterile aesthetic (white backgrounds, black text) that conveys "This is a serious, official tool."
3. **Action-Oriented:** Every element exists to guide the user to the next step (Scan -> Confirm -> Success).

---

## 🖌️ 1. Color Palette (Oklch System)
Do not use complex gradients or heavy shadows. Stick to flat, neutral colors with specific semantic meanings.

| Token | Value (Oklch) | Usage |
| :--- | :--- | :--- |
| **Background** | `oklch(1 0 0)` | Pure white. Main backdrop for cards and the body. |
| **Card Background** | `oklch(1 0 0)` | White. Used inside cards. |
| **Zinc-50 (Light Gray)** | `oklch(0.967 0.001 286.375)` | Used for subtle UI states: "Day" badges, "Free" badges, and input backgrounds. |
| **Zinc-200 (Borders)** | `oklch(0.92 0.004 286.32)` | Thin, 1px borders for inputs and cards. Defines the structure without drawing attention. |
| **Foreground (Text)** | `oklch(0.141 0.005 285.823)` | Near-black (#141414). Do not use pure black. |
| **Muted Text** | `oklch(0.552 0.016 285.938)` | Zinc-500. Used for labels, placeholders, and helper text to create a visual hierarchy. |
| **Primary (Black)** | `oklch(0.21 0.006 285.885)` | The standard "Button" color. High contrast against the white background. |
| **Success (Green)** | `oklch(0.769 0.188 70.08)` | Used strictly for "Check-in Confirmed" states and "Free Event" badges. |
| **Error (Red)** | `oklch(0.577 0.245 27.325)` | Used for "Registration Failed" states and invalid QR codes. |

---

## 🔠 2. Typography
- **Font Family:** `Inter` (Sans-serif) via Geist / Geist Mono.
- **Headings:** 
  - Use `text-zinc-900` for high-impact text (Titles, Day numbers).
  - H1 (Main Title): `text-3xl font-bold`
  - H2 (Card Sections): `text-lg font-semibold`
- **Body Text:** 
  - Use `text-zinc-600` for standard paragraphs.
  - Use `text-sm` for input labels and secondary metadata.
  - Use `text-xs` for timestamps and fine print.
- **Monospace:** Use `font-mono` *only* for technical data (like raw invoice strings or debug info). Never use it for UI labels.

---

## 📐 3. Spacing & Layout
The UI relies entirely on whitespace rather than heavy borders to separate content.

- **Card Structure:**
  - **Border:** `border rounded-xl` (12px radius). Soft corners.
  - **Padding:** Internal padding is always `p-6` or `p-8`. Never less than `p-4`.
  - **Shadow:** Use `shadow-sm` for a subtle lift, but only if the card is overlapping another element. Projector screens do not need heavy shadows.
- **Gaps:** 
  - Use `gap-6` between main sections in a flex/col layout.
  - Use `gap-4` between Form Inputs and Buttons.
  - Use `gap-2` between labels and inputs.
- **Vertical Rhythm:** Always center-align content in the layout. Add `py-4` or `py-6` to allow elements to breathe.

---

## 🧩 4. Core Components

### 4.1 Buttons
- **Primary Action:**
  - `bg-black text-white font-medium rounded-lg py-2.5 px-4`
  - Hover: `bg-zinc-800`
  - Disabled: `opacity-50 cursor-not-allowed`
- **Secondary / Ghost Action:**
  - `border text-zinc-600 font-medium rounded-lg py-2.5 px-4`
  - Hover: `bg-zinc-50`
- **Width:** On mobile/attendee screens, buttons must be `w-full` (full width) to be easily tapped with a thumb.
- **Height:** Maintain a minimum height of `h-10` or `py-2.5` to ensure a large tap target.

### 4.2 Inputs
- **Style:** `border border-zinc-300 rounded-lg px-4 py-2.5 w-full`
- **Focus State:** `focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent`
- **Labels:** Always placed *outside* the input, above it. `block text-sm font-medium text-zinc-700 mb-1`.

### 4.3 Badges (Pills)
Used for status indicators (Day, Free/Paid).
- **Style:** `inline-block bg-zinc-100 px-3 py-1 rounded-full text-xs font-medium text-zinc-600`
- **Semantic Badges:**
  - **Free Event:** `bg-green-50 text-green-700 border border-green-200`
  - **Paid Event:** `bg-blue-50 text-blue-700 border border-blue-200`

### 4.4 Cards & Containers
- **Main Card:** `bg-white p-6 rounded-xl shadow-sm border border-zinc-200`
- **List Item (Check-in list):** 
  - Background: `bg-zinc-50` 
  - Border: `border border-zinc-100`
  - Padding: `py-2 px-3`
  - Layout: `flex justify-between items-center`

---

## 📱 5. Specific Screen Layout Rules

### Screen A: The Organizer (M's Screen)
- **Usage:** Displayed on a laptop, projector, or M's phone. **Must be scannable from 3 feet away.**
- **QR Code Size:** The QR must be at least `280px` or `w-72` in Tailwind.
- **QR Styling:** Wrap in `border-4 border-zinc-200 p-4 rounded-xl bg-white` to define its boundaries clearly.
- **Auto-Refresh:** All data (attendee count, list) must auto-refresh via a `setInterval` (e.g., every 5 seconds). Display a small badge "Auto-refreshes every 5s" to reassure M.

### Screen B: The Attendee (Student's Phone)
- **Usage:** A mobile-first interface (max-width: `max-w-md`).
- **Step-Based Flow:** Do not use modals. Hide/show DIVs based on `step` state (Register -> Pay -> Confirmed).
- **Payment Step:** The QR code for the Lightning Invoice must be pure white background (`bg-white`), centered, and must have a fallback text (e.g., "QR Code Placeholder") if the LND payment request fails to load.
- **Confirmation Step:** Use a massive green checkmark `text-6xl` followed by a clear success message. This is the "reward" for the user.

---

## ⚡ 6. Interaction Patterns (UX Rules)

1. **Immediate Feedback:** Every button click must trigger a `toast` notification (via `sonner`). "Error: Payment failed" or "Success: Checked in!".
2. **Input Validation:** Never allow an empty "Name" or "Email" field to be submitted. 
3. **State Resilience:** If the QR code is scanned and returns an invalid or expired day, show a specific "Invalid Session" card rather than crashing the app.
4. **Touch Targets:** On attendee screens, all buttons must be `w-full` and at least `48px` tall to comply with mobile accessibility standards.

---

## 📦 7. Imports & Utilities (Next.js / Shadcn specific)
- **Preferred UI Library:** shadcn/ui (Radix Primitives) + Tailwind CSS.
- **Icon Library:** `lucide-react`. Use icons sparingly, only to guide the user (e.g., a QR icon, an X to close).
- **State Management:** Use React `useState` and `useEffect` exclusively. No external state libraries (Zustand/Redux) are required for this scale.
- **Routing:** Use Next.js `useSearchParams` to read the `day` query parameter from the URL.