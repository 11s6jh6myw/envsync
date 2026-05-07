// Package watch polls a .env file for content changes and emits
// structured Events whenever the file is modified.
//
// # Overview
//
// A Watcher computes a SHA-256 hash of the target file on each poll tick.
// When the hash differs from the previously recorded value the file is
// re-parsed and an Event is sent on the channel returned by Start.
//
// # Usage
//
//	opts := watch.DefaultOptions()
//	opts.PollInterval = 5 * time.Second
//
//	w := watch.New(".env", opts)
//	ch, err := w.Start()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer w.Stop()
//
//	for ev := range ch {
//		fmt.Printf("%s changed at %s\n", ev.Path, ev.ChangedAt)
//	}
//
// # Notes
//
// Watch uses polling rather than OS-level inotify/FSEvents so that it
// works uniformly across platforms and inside containers.  For low-latency
// use-cases reduce PollInterval; the default is 2 s.
package watch
