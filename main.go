package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// Configurable output directory
var outputDir string

func main() {
	intro()

	// Set default output directory
	outputDir = getOutputDir()
	color.Cyan("[i] Output will be saved to: %s\n", outputDir)

	// Initialize report writer
	var report strings.Builder

	bar := progressbar.NewOptions(6,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[cyan][1/5][reset] Running SUID file scan..."),
	)

	bar.Describe("[cyan][1/6][reset] Running SUID file scan...")
	time.Sleep(1 * time.Second)
	checkSUIDFiles(&report, bar)

	bar.Describe("[cyan][2/6][reset] Checking user accounts...")
	time.Sleep(1 * time.Second)
	checkUsers(&report, bar)

	bar.Describe("[cyan][3/6][reset] Scanning open ports...")
	time.Sleep(1 * time.Second)
	checkOpenPorts(&report, bar)

	bar.Describe("[cyan][4/6][reset] Checking world-writable directories...")
	time.Sleep(1 * time.Second)
	checkWorldWritableDirs(&report, bar)

	bar.Describe("[cyan][5/6][reset] Checking crontabs...")
	time.Sleep(1 * time.Second)
	checkCrontabs(&report, bar)

	bar.Describe("[cyan][6/6][reset] Writing report file...")
	time.Sleep(1 * time.Second)

	bar.Finish()
	color.Green("\n[+] Audit completed. Report saved in %s/audit_report.txt\n", outputDir)

	// Write report to a single file
	writeReport("audit_report.txt", report.String())
}

func intro() {
	logo := `
   ███████╗██╗  ██╗███████╗██╗     █████╗ 
   ██╔════╝██║  ██║██╔════╝██║    ██╔══██╗
   ███████╗███████║█████╗  ██║    ███████║
   ╚════██║██╔══██║██╔══╝  ██║    ██╔══██║
   ███████║██║  ██║███████╗███████╗██║  ██║
   ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝

        SHELA - SHELA Helps Evaluate Linux Access
`
	color.Set(color.FgHiMagenta)
	fmt.Println(logo)
	color.Unset()
}

func getOutputDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func writeReport(filename, content string) {
	path := filepath.Join(outputDir, filename)
	f, err := os.Create(path)
	if err != nil {
		color.Red("[!] Failed to write report: %v", err)
		return
	}
	defer f.Close()
	f.WriteString(content)
}

func checkSUIDFiles(report *strings.Builder, bar *progressbar.ProgressBar) {
	cmd := exec.Command("find", "/", "-perm", "-4000", "-type", "f")
	output, err := cmd.CombinedOutput()
	bar.Add(1)

	files := strings.Split(string(output), "\n")
	count := 0
	for _, f := range files {
		if f != "" {
			count++
		}
	}

	report.WriteString("=== [SUID Files Found] ===\n")
	for _, f := range files {
		if f != "" {
			report.WriteString(f + "\n")
		}
	}
	report.WriteString(fmt.Sprintf("\n[!] Found %d SUID files.\n    -> Recommendation: Review and remove unnecessary SUID bits.\n\n", count))

	if err != nil && count > 0 {
		// Do nothing: we already handled output.
	} else if err != nil {
		color.Red("[!] Error scanning SUID files: %v", err)
	}
}

func checkUsers(report *strings.Builder, bar *progressbar.ProgressBar) {
	file, err := os.Open("/etc/passwd")
	bar.Add(1)
	if err != nil {
		color.Red("[!] Cannot open /etc/passwd")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	report.WriteString("=== [User Accounts Audit] ===\n")
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) > 6 && parts[6] == "/bin/bash" {
			report.WriteString(fmt.Sprintf("User %s can login via shell\n", parts[0]))
		}
	}
	report.WriteString("\n[!] Users with /bin/bash detected.\n    -> Recommendation: Disable shell access for non-admin users.\n\n")
}

func checkOpenPorts(report *strings.Builder, bar *progressbar.ProgressBar) {
	cmd := exec.Command("ss", "-tuln")
	output, err := cmd.CombinedOutput()
	bar.Add(1)
	if err != nil {
		color.Red("[!] Failed to check open ports")
		return
	}
	report.WriteString("=== [Open Ports] ===\n")
	report.WriteString(string(output))
	report.WriteString("\n[!] Open ports detected.\n    -> Recommendation: Close unnecessary ports and use firewall.\n\n")
}

func checkWorldWritableDirs(report *strings.Builder, bar *progressbar.ProgressBar) {
	cmd := exec.Command("sudo", "find", "/", "-type", "d", "-perm", "-0002", "-not", "-path", "/proc/*", "-not", "-path", "/sys/*", "-not", "-path", "/dev/*", "-print")
	output, err := cmd.CombinedOutput() // Mengambil output dan error sekaligus

	if err != nil {
		// Menangani error, khususnya jika error terkait dengan Permission Denied
		if !strings.Contains(string(output), "Permission denied") {
			color.Red("[!] Error finding world-writable dirs: %v", err)
			color.Red("[!] Output: %s", string(output))
		}
	}

	if len(output) == 0 {
		color.Yellow("[!] No world-writable directories found.")
		return
	}

	// Simpan output yang valid
	report.WriteString("=== [World-Writable Directories] ===\n")
	report.WriteString(string(output))
	report.WriteString("\n[!] World-writable directories found.\n    -> Recommendation: Restrict permissions on sensitive directories.\n\n")

}

func checkCrontabs(report *strings.Builder, bar *progressbar.ProgressBar) {
	cmd := exec.Command("ls", "/var/spool/cron")
	output, err := cmd.CombinedOutput()
	bar.Add(1)
	if err != nil {
		color.Red("[!] Crontab check failed")
		return
	}
	report.WriteString("=== [Crontabs] ===\n")
	report.WriteString(string(output))
	report.WriteString("\n[!] User crontabs detected.\n    -> Recommendation: Review periodic tasks for malicious jobs.\n\n")
}
