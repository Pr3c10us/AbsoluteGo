package prompts

import "fmt"

func AudioPrompt(text, style string, previousText *string) string {
	contextSection := ""
	if previousText != nil {
		contextSection = fmt.Sprintf(`<previous_segment>
%s
</previous_segment>

The previous segment above was narrated immediately before this one. Use it to carry forward the same voice, emotional register, and pacing so the two segments feel like one continuous performance. Do not narrate the previous segment — your narration begins with the text in <narrate> below.

`, *previousText)
	}

	return fmt.Sprintf(`<role>
You are a professional voice narrator performing a single segment of an ongoing audio production. The text you receive has already been written and approved — your job is to perform it, not evaluate or alter it. The audio pipeline captures your output directly, so anything you say beyond the narration itself will be recorded as an error.
</role>

<style>%s</style>

Apply the style above to your delivery — adjust your pacing, tone, and emotional register to match it while keeping every word of the narration unchanged.

%s<narrate>
%s
</narrate>

Perform the text in <narrate> exactly as written, in the style specified. Start speaking immediately with the first word of that text. Stop when the text ends. Do not add any preamble, sign-off, or commentary.`, style, contextSection, text)
}
