package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func printHelp() {
	fmt.Println("USAGE: nag time text")
	fmt.Println(`    time: Delay before the reminder (eg "2h", "25m", "1h30m", "48h")`)
	fmt.Println("    text: Reminder text (quote it if it contains spaces)")
}

func scheduleReminder(delay time.Duration, text string) error {
	notifyPath, err := exec.LookPath("notify-send")
	if err != nil {
		return fmt.Errorf("notify-send is missing: %w", err)
	}

	// Convert to seconds rather than passing Go's duration format directly.
	timerDelay := fmt.Sprintf("%.9fs", delay.Seconds())

	cmd := exec.Command(
		"systemd-run",
		"--user",
		"--on-active="+timerDelay,
		"--timer-property=AccuracySec=1s",
		"--collect",
		"--expand-environment=no",
		"--",
		notifyPath,
		"--",
		"nag",
		text,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"could not schedule reminder: %w\n%s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "err: expected exactly 2 arguments")
		printHelp()
		os.Exit(1)
	}

	delay, err := time.ParseDuration(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: invalid time: %v\n", err)
		printHelp()
		os.Exit(1)
	}

	if delay <= 0 {
		fmt.Fprintln(os.Stderr, "err: time must be greater than zero")
		os.Exit(1)
	}

	text := os.Args[2]

	if strings.TrimSpace(text) == "" {
		fmt.Fprintln(os.Stderr, "err: reminder text must not be empty")
		os.Exit(1)
	}

	if err := scheduleReminder(delay, text); err != nil {
		fmt.Fprintln(os.Stderr, "err:", err)
		os.Exit(1)
	}

	fmt.Printf("Reminder scheduled in %s: %s\n", delay, text)
}
