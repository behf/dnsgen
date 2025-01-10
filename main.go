package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	dnsgen "github.com/behf/dnsgen/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// setupLogger configures the logrus logger
func setupLogger(verbose bool) *logrus.Logger {
	logger := logrus.New()
	logger.Out = os.Stdout

	// Set the log level based on the verbose flag
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.WarnLevel)
	}

	// Customize the log format
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger
}

// validateWordLen validates the word length parameter.
func validateWordLen(cmd *cobra.Command, args []string, wordLen int) error {
	if wordLen < 1 || wordLen > 100 {
		return fmt.Errorf("word length must be between 1 and 100")
	}
	return nil
}

// setupGenerator sets up and configures the domain generator.
func setupGenerator(wordlistPath string, wordLen int, logger *logrus.Logger) (*dnsgen.DomainGenerator, error) {
	generator, err := dnsgen.NewDomainGenerator(wordlistPath, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize generator: %w", err)
	}
	logger.Info("Generator initialized successfully")
	return generator, nil
}

// processDomains processes input domains and generates variations.
func processDomains(domains []string, generator *dnsgen.DomainGenerator, wordLen int, fast bool, logger *logrus.Logger) ([]string, error) {
	logger.Info("Generating domain variations...")

	var generated []string
	var wg sync.WaitGroup
	results := make(chan []string)

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			variations := generator.Generate([]string{domain}, wordLen, fast)
			logger.Debugf("Generated variations for domain %s: %v", domain, variations)
			results <- variations
		}(domain)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	uniqueDomains := make(map[string]bool)
	for result := range results {
		for _, domain := range result {
			uniqueDomains[domain] = true
		}
	}

	for domain := range uniqueDomains {
		generated = append(generated, domain)
	}

	logger.Info("Finished generating domain variations")
	return generated, nil
}

// writeOutput writes generated domains to output file or stdout.
func writeOutput(domains []string, outputPath string, logger *logrus.Logger) error {
	sort.Strings(domains)

	if outputPath != "" {
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, domain := range domains {
			fmt.Fprintln(writer, domain)
		}
		writer.Flush()
		logger.Info("Results written to", outputPath)
	} else {
		for _, domain := range domains {
			fmt.Println(domain)
		}
	}
	return nil
}

func main() {
	var wordLen int
	var wordlistPath string
	var fast bool
	var outputPath string
	var verbose bool

	var rootCmd = &cobra.Command{
		Use:   "dnsgen",
		Short: "Generate DNS name permutations for domain discovery.",
		Args:  cobra.ArbitraryArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return validateWordLen(cmd, args, wordLen)
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := setupLogger(verbose)

			var inputDomains []string
			var err error

			// Check if input is coming from a pipe
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				// Data is being piped, read from stdin
				reader := bufio.NewReader(os.Stdin)
				for {
					line, readErr := reader.ReadString('\n')
					if readErr != nil {
						if readErr == io.EOF {
							break // End of input
						}
						logger.WithError(readErr).Error("Error reading from stdin")
						os.Exit(1)
					}
					line = strings.TrimSpace(line)
					if line != "" {
						inputDomains = append(inputDomains, line)
					}
				}
			} else {
				// No data is piped, expect a filename as an argument
				if len(args) < 1 {
					logger.Error("Missing input filename")
					cmd.Usage()
					os.Exit(1)
				}
				inputFilePath := args[0]
				inputFile, err := os.Open(inputFilePath)
				if err != nil {
					logger.WithError(err).Error("Error opening input file")
					os.Exit(1)
				}
				defer inputFile.Close()

				scanner := bufio.NewScanner(inputFile)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())
					if line != "" {
						inputDomains = append(inputDomains, line)
					}
				}
				if err := scanner.Err(); err != nil {
					logger.WithError(err).Error("Error reading input file")
					os.Exit(1)
				}
			}

			logger.Infof("Read %d domains from input", len(inputDomains))

			generator, err := setupGenerator(wordlistPath, wordLen, logger)
			if err != nil {
				logger.WithError(err).Error("Failed to setup generator")
				os.Exit(1)
			}

			// Register the default permutators
			generator.RegisterDefaultPermutators()

			generatedDomains, err := processDomains(inputDomains, generator, wordLen, fast, logger)
			if err != nil {
				logger.WithError(err).Error("Failed to process domains")
				os.Exit(1)
			}

			logger.Infof("Generated %d unique domain variations", len(generatedDomains))
			err = writeOutput(generatedDomains, outputPath, logger)
			if err != nil {
				logger.WithError(err).Error("Failed to write output")
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().IntVarP(&wordLen, "wordlen", "l", 6, "Minimum length of custom words extracted from domains.")
	rootCmd.PersistentFlags().StringVarP(&wordlistPath, "wordlist", "w", "", "Path to custom wordlist file.")
	rootCmd.PersistentFlags().BoolVarP(&fast, "fast", "f", false, "Use fast generation mode (fewer permutations).")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "Output file path.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging.")

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Failed to execute command")
	}
}
