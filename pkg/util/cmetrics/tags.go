package cmetrics

import (
	"os"
	"strings"
)

var shortHostName = ""

// GetShortHostname gives the hostname for use with tags. It's cached so we
// only ask the kernel once.
func GetShortHostname() string {
	if shortHostName == "" {
		shortHostName, _ = os.Hostname()
		// shortHostName = strings.Replace(shortHostName, ".", "_", -1) // TODO: someone wanted this but didn't finish implementing it
	}
	return shortHostName
}

// TagsMap takes an array of tags and turns it into a map.
func TagsMap(tags []string) map[string]string {
	tagsMap := make(map[string]string)
	for _, t := range tags {
		// TODO(stefan): error out on bad tag names, or fix them
		// http://opentsdb.net/docs/build/html/user_guide/writing/index.html
		// Only the following characters are allowed: a to z, A to Z, 0 to 9, -, _, ., / or Unicode letters (as per the specification)
		pts := strings.SplitN(t, "=", 2)
		if len(pts) > 1 {
			tagsMap[pts[0]] = pts[1]
		}
	}
	return tagsMap
}

// ExpandTags takes a metric name with optional tags in the name, combines it with other tags,
// adds default tags, etc. You could probably simplify this a lot.
func ExpandTags(metricNameBase, metricNamePrefix, hostTag string, tagsBase, tagsExtra map[string]string) (string, map[string]string) {
	name := metricNameBase
	tags := make(map[string]string)

	for k, v := range tagsBase {
		tags[k] = v
	}

	tags["host"] = hostTag

	for k, v := range tagsExtra {
		tags[k] = v
	}

	nPts := strings.Split(name, "^")
	if len(nPts) > 1 {
		name = nPts[0]
		for _, np := range nPts[1:] {
			npr := strings.Split(np, "=")
			// e.g. the name was "mymetric^$CID=1234" and
			// the tags array should have "cid=$CID".
			// Don't use this. Prefer the branch below.
			if npr[0][0] == '$' {
				for k, v := range tags {
					if v == npr[0] {
						tags[k] = npr[1]
					}
				}
			} else {
				// e.g. the name was "mymetric^sometag=somevalue".
				// We will send a tag sometag=somevalue.
				tags[npr[0]] = npr[1]
			}
		}
	}

	if metricNamePrefix != "" {
		name = metricNamePrefix + "." + name
	}

	return name, tags
}
