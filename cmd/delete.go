/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a particular oktaasa host id",
	Long: `Searches project for a hostname that matches, then deletes that ID

Very brute force.`,
	Run: deleteRun,
}

var hostname string
var project string

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")
	deleteCmd.PersistentFlags().StringVar(&project, "project", "", "scaelft project name")
	deleteCmd.PersistentFlags().StringVar(&hostname, "hostname", "", "scaelft hostname")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func deleteRun(cmd *cobra.Command, args []string) {
	bearer, err := getToken(viper.GetString("OKTAASA_KEY"), viper.GetString("OKTAASA_KEY_SECRET"), viper.GetString("OKTAASA_TEAM"))
	if err == nil {
		//fmt.Println("Bearer token %s", bearer)
	} else {
		fmt.Println("Error getting token #v", err)
	}
	//fmt.Println("project ", project, " hostname ", hostname)
	list, err := get_servers(bearer, viper.GetString("OKTAASA_TEAM"), project)
	if err != nil {
		//fmt.Println(fmt.Errorf("Error getting server list. error:%v", err))
		return
	}

	ids := get_ids_for_hostname(hostname, list)

	if len(ids) == 0 {
		//      return fmt.Errorf("Error, ScaleFT api returned no servers that matched hostname:%s", hostname)
		//      This should not happen, but if it does, it's ok?
		log.Printf("[WARN] No servers matched for Hostname:%s", hostname)
		return
	}

	for _, id := range ids {
		if id != "" {
			err := delete_server(bearer, viper.GetString("OKTAASA_TEAM"), project, id)
			//log.Printf("looped %s", id)
			if err != nil {
				//	log.Printf("[WARN] Failed to delete server with hostname: %s at ScaleFT ID:%s, error:%s", hostname, id, err)
				//              return fmt.Errorf("Error deleting server at id:%s and key_team:%s project: %s error:%v", id, key_team, project, err)
			}
		}
	}

	return
}

const api string = "https://app.scaleft.com/v1/teams/"

type Body struct {
	Key_id     string `json:"key_id"`
	Key_secret string `json:"key_secret"`
}

type Bearer struct {
	Bearer_token string `json:"bearer_token"`
}

func getToken(key_id string, key_secret string, key_team string) (string, error) {
	p := &Body{key_id, key_secret}
	jsonStr, err := json.Marshal(p)
	url := api + key_team + "/service_token"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "error", fmt.Errorf("Error getting token key_id:%s key_team:%s status:%s error:%v", key_id, key_team, string(resp.Status), err)
	}

	defer resp.Body.Close()
	b := Bearer{}
	json.NewDecoder(resp.Body).Decode(&b)

	return b.Bearer_token, err
}

type Server struct {
	Id              string                 `json:"id"`
	ProjectName     string                 `json:"project_name"`
	Hostname        string                 `json:"hostname"`
	AltNames        []string               `json:"alt_names"`
	AccessAddress   string                 `json:"access_address"`
	OS              string                 `json:"os"`
	RegisteredAt    time.Time              `json:"registered_at"`
	LastSeen        time.Time              `json:"last_seen"`
	CloudProvider   string                 `json:"cloud_provider"`
	SSHHostKeys     []string               `json:"ssh_host_keys"`
	BrokerHostCerts []string               `json:"broker_host_certs"`
	InstanceDetails map[string]interface{} `json:"instance_details"`
	State           string                 `json:"state"`
}

type Servers struct {
	List []*Server `json:"list"`
}

func get_servers(bearer_token string, key_team string, project string) (Servers, error) {
	client := &http.Client{}
	url := api + key_team + "/projects/" + project + "/servers"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return Servers{}, fmt.Errorf("Error listing servers: key_team:%s project:%s status:%s error:%v", key_team, project, string(resp.Status), err)
	}

	s := struct {
		List []*Server `json:"list"`
	}{nil}

	json.NewDecoder(resp.Body).Decode(&s)
	return s, err
}

func delete_server(bearer_token string, key_team string, project string, server_id string) error {
	client := &http.Client{}
	url := api + key_team + "/projects/" + project + "/servers/" + server_id
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error deleting server:%s status:%s error:%v", server_id, string(resp.Status), err)
	}
	//log.Println("delete server URL was %s", url)
	return nil
}

func get_ids_for_hostname(hostname string, server_list Servers) []string {
	filtered := make([]string, len(server_list.List))
	//log.Println("Looking for ", hostname)
	for i, l := range server_list.List {
		//log.Println("checking ", l.Hostname)
		if hostname == l.Hostname {
			//log.Println(hostname, " = ", l.Hostname, ". ID was ", l.Id)
			filtered[i] = l.Id
		}
	}
	return filtered
}
