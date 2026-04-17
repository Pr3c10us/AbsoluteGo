package prompts

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SplitScriptResult struct {
	Chapter int    `json:"chapter"`
	Page    int    `json:"page"`
	Script  string `json:"script"`
	Panel   int    `json:"panel"`
	Effect  string `json:"effect"`
}

func ParseSplitScriptResponse(raw string) ([]SplitScriptResult, error) {
	trimmed := strings.TrimSpace(raw)
	trimmed = stripCodeFences(trimmed)

	var results []SplitScriptResult
	if err := json.Unmarshal([]byte(trimmed), &results); err != nil {
		return nil, fmt.Errorf("failed to parse split script JSON: %w", err)
	}

	if err := validateResults(results); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return results, nil
}

func stripCodeFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "'''") {
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		if strings.HasSuffix(s, "'''") {
			s = s[:len(s)-3]
		}
		s = strings.TrimSpace(s)
	}
	return s
}

func validateResults(results []SplitScriptResult) error {
	if len(results) == 0 {
		return fmt.Errorf("empty results array")
	}

	allEffects := map[string]bool{
		"panLeft": true, "panRight": true,
		"panUp": true, "panDown": true,
		"zoomIn": true, "zoomOut": true,
	}

	for i, r := range results {
		if r.Chapter < 1 {
			return fmt.Errorf("entry %d: invalid chapter number %d", i, r.Chapter)
		}
		if r.Page < 1 {
			return fmt.Errorf("entry %d: invalid page number %d", i, r.Page)
		}
		if r.Panel < 1 {
			return fmt.Errorf("entry %d: invalid panel number %d", i, r.Panel)
		}
		if r.Script == "" {
			return fmt.Errorf("entry %d: empty script segment", i)
		}
		if !allEffects[r.Effect] {
			return fmt.Errorf("entry %d: invalid effect %q", i, r.Effect)
		}
	}

	return nil
}

func SplitScriptPrompt(script string) string {
	return fmt.Sprintf(`<role>
You are a script editor specialising in synchronising audio drama scripts with comic book pages and panels. Your output drives an automated video production pipeline: each JSON entry you produce maps directly to a panel image that will be animated and narrated. Errors in chapter/page numbers, panel assignments, or effect types cause production failures downstream, so precision is essential.
</role>

<task>
You will receive comic book page images and a complete TTS audio drama script. Your job is to split the script into segments, match each segment to the correct chapter and page, assign each segment to a specific panel on that page, and select a motion effect for that panel. The result is a JSON array that the video pipeline consumes directly — your entire response must be that array and nothing else.
</task>

<inputs>
<comic_pages>Attached as images. Each page displays its chapter and page number at the bottom right (e.g. "Chapter 1 Page 2" or "Ch. 1 P. 2"). Panels on each page are labeled with numbers by an automated system.</comic_pages>
<script>
%s
</script>
</inputs>

<workflow>
Work through these four phases in order. Each builds on the previous — do not skip ahead.

<phase_1_page_inventory>
Before reading the script, scan every page image and build a mental inventory:
1. Extract the chapter number and page number from the bottom right of every page. Pages with identical page numbers but different chapter numbers are entirely distinct — treat them as separate entries.
2. Classify each page as one of: cover/title/splash (only eligible for preamble segments), story page with panels (eligible for body segments), or skip page (credits, ads, letters, previews, blank — never used).
3. For each story page, note how many labeled panels it contains.

This inventory prevents chapter/page confusion when you reach the assignment phase.
</phase_1_page_inventory>

<phase_2_panel_reading_order>
For each story page, determine the true visual reading order of its panels before touching the script. The automated labeling system assigns numbers arbitrarily — they do not reflect reading order.

To find true reading order:
1. Identify whether the comic is Western (left-to-right, top-to-bottom) or manga (right-to-left, top-to-bottom). Default to Western unless the art style or text direction signals otherwise.
2. Mentally divide the page into rows by vertical position. Within each row, sort panels by horizontal position according to reading direction.
3. Record the sequence of panel labels in the order they should be read — this is the order you will assign script segments.

<example>
A Western page has three panels. Panel label 2 is top-left, label 1 is top-center, label 3 is top-right. True reading order is: 2 → 1 → 3. You assign script to panels in that sequence, even though the labels are out of numerical order.
</example>

Once you have passed a panel in visual sequence, do not return to it. Moving backward in reading order is not allowed. Repeating a panel consecutively is allowed when a scene needs to linger (extended dialogue, a held tension beat).
</phase_2_panel_reading_order>

<phase_3_script_splitting_and_assignment>
Now read the script and split it into segments, matching each to a chapter, page, and panel.

<preamble_vs_body>
Scripts often open each chapter with preamble — a hook, context bridge, or intro narration that precedes the chapter's panel-by-panel events. Assign preamble segments to the cover, title, or splash page of the matching chapter. Once the narration begins describing actual panel events (characters acting, dialogue occurring), switch to story pages. Never assign body segments to cover or title pages.
</preamble_vs_body>

<splitting_guidelines>
Split at natural breakpoints: chapter transitions, scene or location changes, significant time jumps, perspective shifts between characters, and major dramatic beats. Do not split mid-sentence unless a dramatic pause makes it appropriate.

Each segment should be long enough to accompany a panel — typically one to five lines of narration. When a page has multiple panels, distribute the script across them rather than assigning everything to one panel.

Match each segment to the page whose visual content best represents it: the page showing the characters speaking, the action being described, or the emotional beat being conveyed. Skipping pages that have no narrative equivalent is allowed. Returning to a page for multiple distinct moments is allowed.
</splitting_guidelines>

<panel_shape_and_effect_selection>
For each panel you assign, classify its shape by comparing width to height, then select the effect whose direction points toward the action or focal point in that panel.

Shape determines which effect types are valid:

| Panel shape | Condition | Valid effects |
|---|---|---|
| Horizontal | Width > Height | panLeft, panRight |
| Vertical | Height > Width | panUp, panDown |
| Square | Width ≈ Height | zoomIn, zoomOut |

Direction is determined by the visual content of the panel — not chosen at random. Ask: where is the action, focal character, or point of narrative interest within this panel? The camera should move toward it.

<examples>
<example>
<situation>Horizontal panel. A character delivers a punch to an enemy standing on the left side of the frame.</situation>
<choice>panLeft — the camera moves toward the action on the left.</choice>
</example>

<example>
<situation>Horizontal panel. A car speeds across a highway toward the right edge.</situation>
<choice>panRight — the camera follows the direction of movement.</choice>
</example>

<example>
<situation>Horizontal panel. Two characters face each other in conversation. The script segment covers the left character's dialogue.</situation>
<choice>panLeft — the camera moves toward the speaking character.</choice>
</example>

<example>
<situation>Horizontal panel. Wide establishing shot of a landscape with no dominant focal point.</situation>
<choice>panRight — default for establishing shots, mimicking the natural left-to-right reading sweep.</choice>
</example>

<example>
<situation>Vertical panel. A character's full body is shown in a power pose, head at top, feet at bottom.</situation>
<choice>panDown — the camera scans from head to toe, revealing the character's full presence.</choice>
</example>

<example>
<situation>Vertical panel. A plane ascending into clouds, action concentrated at the top.</situation>
<choice>panUp — the camera follows the upward movement.</choice>
</example>

<example>
<situation>Vertical panel. A character plummeting from a rooftop.</situation>
<choice>panDown — the camera follows the fall.</choice>
</example>

<example>
<situation>Square panel. Close-up of a character's terrified face.</situation>
<choice>zoomIn — intensifies the emotional moment by tightening the frame.</choice>
</example>

<example>
<situation>Square panel. Two armies facing off across a battlefield.</situation>
<choice>zoomOut — pulls back to show scale and context.</choice>
</example>

<example>
<situation>Square panel. A ticking bomb on a table.</situation>
<choice>zoomIn — focuses dramatic attention on the threat.</choice>
</example>

<example>
<situation>Square panel. A character walking away into the distance, end of a scene.</situation>
<choice>zoomOut — creates a sense of departure and closure.</choice>
</example>
</examples>

When the same panel is repeated consecutively, alternate the effect direction to prevent visual monotony:
- Horizontal: panLeft → panRight → panLeft
- Vertical: panUp → panDown → panUp
- Square: zoomIn → zoomOut → zoomIn
</panel_shape_and_effect_selection>
</phase_3_script_splitting_and_assignment>

<phase_4_self_check>
Before producing your final output, verify the following. Revise any entry that fails.

Script completeness:
- Every line of the input script appears in exactly one segment — nothing dropped, nothing duplicated.
- Segments appear in the same chronological order as the source script.

Chapter and page accuracy:
- Every entry's chapter and page values were read from the bottom right of an actual page image — not inferred or assumed.
- No two entries treat Chapter 1 Page 2 and Chapter 2 Page 2 as the same page.
- Preamble segments are assigned only to cover/title/splash pages of the correct chapter.
- Body segments are assigned only to story pages with panels.
- No entry references a skip page (credits, ads, letters, blank, previews).

Panel order:
- Panels on each page are assigned in true visual reading order, not numerical label order.
- No entry goes backward in visual reading sequence on its page.
- If a panel repeats consecutively, effect directions alternate.
- No page has its entire script assigned to a single panel when multiple panels exist.

Effects:
- Every horizontal panel uses panLeft or panRight.
- Every vertical panel uses panUp or panDown.
- Every square panel uses zoomIn or zoomOut.
- Every effect direction was chosen based on where the action or focal point sits in the panel — not randomly.

JSON validity:
- The output is a syntactically valid JSON array with no text outside it.
- All string values preserve the original script formatting, including audio tags and break tags.
- Internal quotes are properly escaped.
</phase_4_self_check>
</workflow>

<output_format>
Produce a JSON array. Each entry has exactly these five fields:

- "chapter": integer — chapter number from the bottom right of the page
- "page": integer — page number within that chapter, from the bottom right
- "script": string — the exact script segment, preserving all [audio tags] and break tags
- "panel": integer — the label number of the panel as shown in the image
- "effect": string — one of: panLeft, panRight, panUp, panDown, zoomIn, zoomOut

Your entire response must be the JSON array. Start with [ and end with ]. Do not include code fences, explanatory text, section headers, or commentary of any kind.

<example>
[
  {
    "chapter": 1,
    "page": 1,
    "script": "[tension] The nightmare begins in red light.",
    "panel": 2,
    "effect": "zoomIn"
  },
  {
    "chapter": 1,
    "page": 1,
    "script": "A boy pounds on a door, screaming for his father,",
    "panel": 1,
    "effect": "panRight"
  },
  {
    "chapter": 1,
    "page": 1,
    "script": "but the man on the other side will never answer.",
    "panel": 3,
    "effect": "panUp"
  },
  {
    "chapter": 1,
    "page": 3,
    "script": "[fast] Bats swarm from the darkness, consuming everything.",
    "panel": 1,
    "effect": "panRight"
  },
  {
    "chapter": 1,
    "page": 3,
    "script": "[short pause] Then he wakes.",
    "panel": 2,
    "effect": "zoomOut"
  },
  {
    "chapter": 2,
    "page": 1,
    "script": "[neutral] A new chapter begins. The sun rises over Gotham.",
    "panel": 1,
    "effect": "panUp"
  }
]
</example>
</output_format>`, script)
}
