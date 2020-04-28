package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	shell "github.com/RTradeLtd/go-ipfs-api"
	"github.com/RTradeLtd/thc"
	"github.com/urfave/cli"
)

var (
	user, pass, url string
	ipfsAPI         = "https://api.ipfs.temporal.cloud"
	ipfsAPIDev      = "https://dev.api.ipfs.temporal.cloud"
	dev             bool
)

func main() {
	app := cli.NewApp()
	app.Flags = loadFlags()
	app.Commands = loadCommands()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func loadFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "user.name",
			Usage:       "the username of your temporal account",
			Destination: &user,
		},
		cli.StringFlag{
			Name:        "user.pass",
			Usage:       "the password of your temporal account",
			Destination: &pass,
		},
		cli.BoolFlag{
			Name:        "dev",
			Usage:       "toggle usage of the dev environment",
			Destination: &dev,
		},
		cli.StringFlag{
			Name:  "custom.url",
			Usage: "specify a custom url to override defaults",
			Value: "",
		},
	}
}

func loadCommands() cli.Commands {
	return []cli.Command{
		{
			Name:  "warp-status",
			Usage: "get warp status",
			Action: func(c *cli.Context) error {
				url := getIPFSAPI(c)
				req, err := http.NewRequest("POST", url+"/ping", nil)
				if err != nil {
					return err
				}
				resp, err := httpDo(req)
				if err != nil {
					return err
				}
				fmt.Println(resp)
				return nil
			},
		},
		{
			Name:  "temporal-status",
			Usage: "get temporal status",
			Action: func(c *cli.Context) error {
				url := getTHCURL(c)
				req, err := http.NewRequest("GET", url+"/systems/check", nil)
				if err != nil {
					return err
				}
				resp, err := httpDo(req)
				if err != nil {
					return err
				}
				fmt.Println(resp)
				return nil
			},
		},
		{
			Name:  "upload",
			Usage: "upload commands",
			Subcommands: cli.Commands{
				{
					Name:        "dir",
					Usage:       "upload directory",
					Description: "uploads a directory and pins for specified duration",
					Action: func(c *cli.Context) error {
						fmt.Println("starting processing")
						v2 := thc.NewV2(user, pass, getTHCURL(c))
						if err := v2.Login(); err != nil {
							fmt.Println("failed to login")
							return err
						}
						fmt.Println("getting jwt")
						jwt, err := v2.GetJWT()
						if err != nil {
							fmt.Println("failed to get jwt")
							return err
						}
						shell := shell.NewDirectShell(getIPFSAPI(c)).WithAuthorization(jwt)
						if c.String("dir") == "" {
							return errors.New("dir flag is empty")
						}
						var hash string
						fmt.Println("uploading directory")
						if hash, err = shell.AddDir(c.String("dir")); err != nil {
							return err
						}
						fmt.Println("hash of directory: ", hash)
						fmt.Println("pinning directory hash")
						if _, err := v2.PinExtend(hash, c.String("hold.time")); err != nil {
							return err
						}
						fmt.Println("successfully pinned directory hash")
						return nil
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the directory to upload",
						},
						cli.StringFlag{
							Name:  "hold.time, ht",
							Usage: "the number of months to pin for",
							Value: "1",
						},
					},
				},
				{
					Name:        "file",
					Usage:       "upload a file",
					Description: "uploads a file and pins for specified duration",
					Action: func(c *cli.Context) error {
						v2 := thc.NewV2(user, pass, getTHCURL(c))
						if err := v2.Login(); err != nil {
							return err
						}
						hash, err := v2.FileAdd(
							c.String("file.name"),
							thc.FileAddOpts{HoldTime: c.String("hold.time")},
						)
						if err != nil {
							return err
						}
						fmt.Println("hash of file", hash)
						fmt.Println(hash + " uploaded and pinned for 1 month")
						return nil
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "file.name, fn",
							Usage: "the file to upload",
						},
						cli.StringFlag{
							Name:  "hold.time, ht",
							Usage: "the number of months to pin for",
							Value: "1",
						},
					},
				},
			},
		},
		{
			Name:        "pin",
			Usage:       "pin a hash",
			Description: "pins a hash for the specified duration",
			Action: func(c *cli.Context) error {
				v2 := thc.NewV2(user, pass, getTHCURL(c))
				if err := v2.Login(); err != nil {
					return err
				}
				if c.String("hash") == "" {
					return errors.New("hash flag cant be empty")
				}
				if _, err := v2.PinAdd(c.String("hash"), c.String("hold.time")); err != nil {
					return err
				}
				fmt.Println(c.String("hash") + " pinned for 1 month")
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hash",
					Usage: "the hash to pin",
				},
				cli.StringFlag{
					Name:  "hold.time, ht",
					Usage: "the number of months to pin for",
					Value: "1",
				},
			},
		},
		{
			Name:  "lens",
			Usage: "use the lens search engine",
			Subcommands: cli.Commands{
				{
					Name:  "index",
					Usage: "index a hash",
					Action: func(c *cli.Context) error {
						v2 := thc.NewV2(user, pass, getTHCURL(c))
						if err := v2.Login(); err != nil {
							return err
						}
						_, err := v2.IndexHash(c.String("hash"), c.Bool("reindex"))
						return err
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "hash",
							Usage: "the hash to index",
						},
						cli.BoolFlag{
							Name:  "reindex",
							Usage: "force a reindex",
						},
					},
				},
				{
					Name:  "search",
					Usage: "search the lens index",
					Action: func(c *cli.Context) error {
						v2 := thc.NewV2(user, pass, getTHCURL(c))
						if err := v2.Login(); err != nil {
							return err
						}
						resp, err := v2.SearchLens(c.String("query"))
						if err != nil {
							return err
						}
						fmt.Printf("results\n%+v\n", resp.Response.Results)
						return nil
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "query",
							Usage: "the query to perform",
						},
					},
				},
			},
		},
	}
}

// readFile is used to read the content of the file into an array
func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getTHCURL(c *cli.Context) string {
	if c.GlobalString("custom.url") != "" {
		return c.GlobalString("custom.url")
	}
	if dev {
		return thc.DevURL
	}
	return thc.ProdURL
}

func getIPFSAPI(c *cli.Context) string {
	if c.GlobalString("custom.url") != "" {
		return c.GlobalString("custom.url")
	}
	if dev {
		return ipfsAPIDev
	}
	return ipfsAPI
}

func httpDo(req *http.Request) (string, error) {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
