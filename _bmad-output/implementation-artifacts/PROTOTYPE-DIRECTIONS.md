# Interactive Prototype Directions — ASD v1.0

**Status:** Ready for Implementation  
**Date:** 31 janvier 2026  
**Duration:** Week 1 (5 days)  
**Tools:** Figma or Framer (recommended Framer for interactive)

---

## Overview

Create **5 interactive prototypes** from UX specification, validating all user journeys end-to-end with real interaction patterns.

**Deliverables:**
- 5 clickable journeys (find & download, browse trending, set & automate, monitor & debug, manage subscriptions)
- Mode Simple/Expert toggle working
- All key interactions: buttons, forms, modals, feedback
- User testing ready (share link with 2-3 test users)

---

## Journey Flows to Prototype

### Journey 1: Find & Download Quick (Alex)
**Screens:**
1. Home/Search tab
2. Search results (anime list)
3. Anime detail (title, synopsis, [Download] button)
4. "When?" modal (Now | Scheduled | Rule)
5. Download progress (%, speed, ETA)
6. Success confirmation

**Interactions:**
- Search → Enter anime name → Click anime → Click [Download] → Select "Now" → Download starts
- Progress visible, completion notification

**Mode:** Simple (no quality dialog, just download)

---

### Journey 2: Browse Trending (Alex)
**Screens:**
1. Home tab (default view)
2. Trending section (scrollable cards)
3. New Releases section
4. Popular Series section
5. Anime detail (from card click)
6. Action choice ([Download] or [Subscribe])

**Interactions:**
- Scroll trending → Click anime card → Read details → Click [Download] → Journey 1 or [Subscribe] → Journey 3

---

### Journey 3: Set & Automate (Maya)
**Screens:**
1. Automations tab
2. My Rules list (empty or with existing rules)
3. [Create Rule] button flow:
   - Step 1: Trigger type (AniList or Manual list)
   - Step 2: Quality selection (1080p, 720p, etc.)
   - Step 3: Language selection (JP, FR, etc.)
   - Step 4: Timing (on air date, custom, etc.)
4. Review summary
5. Rule created confirmation

**Interactions:**
- Click [Create Rule] → Step 1 form → [Continue] → Step 2 → [Continue] → Step 3 → [Continue] → Review → [Create] → Success toast

**Validation:**
- Real-time validation (at least one quality, one language, valid timing)
- Error messages inline (red text below field)

---

### Journey 4: Monitor & Debug (Maya)
**Screens:**
1. Dashboard tab
2. My Automations list with status badges (Active ✓, Alert 🚨, Running 🔄)
3. Rule details (click a rule)
4. Execution logs (clickable, shows full log)
5. Edit rule inline (click [Edit] button)
6. Save rule changes

**Interactions:**
- Click rule with Alert status → See failure reason → Click [View Logs] → Full log visible → Click [Edit] → Adjust quality/language → [Save] → Confirmation

**Visual:**
- Status badges: Icon + label (not color alone)
- Logs: Monospace font, timestamp + message, syntax highlighted

---

### Journey 5: Manage Subscriptions (Maya)
**Screens:**
1. My Rules list
2. [Edit] action on rule card
3. Edit form (quality, language, timing)
4. [Save] or [Delete] buttons
5. Confirmation modals for destructive actions

**Interactions:**
- Hover rule card → [Edit], [Delete], [•••] visible → Click [Delete] → "Are you sure?" modal → [Delete] confirms

---

## Mode Toggle Feature

**Implementation:**
- Toggle switch: Top-right corner, always visible
- States: "Simple" ↔ "Expert"
- Behavior:
  - Simple: Hide logs, complexity, advanced options
  - Expert: Show all details, logs, advanced settings
  - Smooth transition (fade 200ms)

**Example Differences:**
- Journey 1 Simple: One [Download] button, download starts
- Journey 1 Expert: [When?] modal, quality preview, log access

---

## Visual Direction: "Harmony"

**Color Scheme:**
- Background: Dark #0A0E1A
- Text: Light #F5F5F5
- Primary Button: Magenta #D946EF
- Secondary Button: Cyan #06B6D4
- Error: Red #F87171
- Success: Green #4ADE80
- Warning: Amber #FBBF24

**Typography:**
- Headings: Noto Serif JP (bold)
- Body: Inter 14px
- Monospace (logs): SF Mono
- Buttons: Inter semibold (500)

**Spacing:**
- Base unit: 4px
- Gaps: 16px (mobile), 20px (tablet), 24px (desktop)
- Padding: 12px (buttons), 16px (cards)

---

## Component Library (from UX spec)

Use **shadcn/ui** as base, customize with Sakura Night tokens:

**Foundation:**
- Button (primary: magenta, secondary: cyan)
- Input, Select, Form
- Card, Badge, Modal, Tabs
- Progress bar, Toast

**Custom (specified in UX doc):**
1. **StatusBadge** — ✓ Active / 🚨 Alert / 🔄 Running
2. **RuleCard** — Rule with inline actions
3. **FormStepper** — Multi-step rule creation
4. **LogViewer** — Full execution logs
5. **ModeToggle** — Simple ↔ Expert switch
6. **DownloadProgress** — Download tracking with ETA

---

## Prototype Validation Checklist

**Before User Testing:**

- [ ] All 5 journeys clickable end-to-end
- [ ] Mode toggle switches all screens properly
- [ ] Forms validate (at least one field required)
- [ ] Buttons show feedback (hover, click states)
- [ ] Modals open/close correctly
- [ ] Progress bars animate
- [ ] Status badges display correctly
- [ ] Mobile (375px) responsive view works
- [ ] All text readable (no overflow)
- [ ] Links work (internal navigation)

---

## User Testing Script (Week 1, End)

**Recruit 2-3 test users:**
- 1 like "Alex" (casual, just wants to download)
- 1-2 like "Maya" (expert, power user)

**Test Scenarios:**

**Alex Testing:**
1. "You want to download Demon Slayer episode 12. Can you do it?"
   → Expected: Find anime → Click download → Done (Simple mode)

2. "Browse for something new to watch"
   → Expected: Use trending section, discover anime

**Maya Testing:**
1. "Set up automatic downloads for any new anime in your AniList"
   → Expected: Create rule (trigger, quality, timing)

2. "One rule failed. Can you see what went wrong and fix it?"
   → Expected: View logs, understand error, edit rule

**Feedback to Collect:**
- Was it obvious what to do next?
- Did the interface feel responsive/fast?
- Any confusing parts?
- What would make it better?

---

## Next Steps

**If Testing Passes:**
→ Move to React implementation immediately (Week 2)

**If Issues Found:**
→ Quick refinements to prototype
→ Re-test with same users
→ Then implement

---

## Timeline

```
Monday (Day 1):    Setup Figma/Framer, create Home + Search screens
Tuesday (Day 2):   Anime detail, "When?" modal, download progress
Wednesday (Day 3): Automations tab, rule creation flow
Thursday (Day 4):  Dashboard, rule details, logs viewer, edit flow
Friday (Day 5):    Polish interactions, mode toggle, user testing prep

Week 2+: React implementation based on validated prototype
```

---

**Ready to start prototyping? Figma link or Framer project to create?** 🎨

