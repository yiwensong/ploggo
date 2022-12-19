package main

import (
	fmt "fmt"
	os "os"
	path "path"

	errors "github.com/pkg/errors"
	cobra "github.com/spf13/cobra"
	avalon "github.com/yiwensong/ploggo/avalon"
	storage "github.com/yiwensong/ploggo/avalon/storage"
)

var rootCmd = &cobra.Command{
	Use:   "avalon",
	Short: "record an avalon game",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmdExec(cmd, args)
	},
}

var minions *[]string
var loyalServants *[]string
var winnerString *string
var isTest *bool

func init() {
	minions = rootCmd.Flags().StringArrayP(
		"minions-of-mordred",
		"i",
		[]string{},
		"usernames of the non-roled red",
	)
	loyalServants = rootCmd.Flags().StringArrayP(
		"loyal-servants",
		"l",
		[]string{},
		"usernames of non-roled blue",
	)
	winnerString = rootCmd.Flags().StringP(
		"winner",
		"w",
		"",
		"the winner of the game. must be 'blue' or 'red'",
	)
	isTest = rootCmd.Flags().BoolP(
		"test",
		"",
		false,
		"use this flag to enable testing mode",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf(err.Error())
	}
}

var BASE_PATH = path.Join(os.Getenv("HOME"), ".avalon")
var TEST_PATH = path.Join(os.Getenv("HOME"), ".avalon_test")

func rootCmdExec(cmd *cobra.Command, args []string) error {
	// merlin := flags.StringP("merlin", "m", "", "username of merlin")
	// percival := flags.StringP("percival", "p", "", "username of percival")
	// morgana := flags.StringP("morgana", "r", "", "username of morgana")

	// mordred := flags.StringP("mordred", "d", "", "username of mordred")
	// assassin := flags.StringP("assassin", "a", "", "username of the assassin")
	// oberon := flags.StringP("oberon", "o", "", "username of oberon")

	gameWinner := avalon.NoTeam
	if *winnerString == "blue" {
		gameWinner = avalon.Blue
	} else if *winnerString == "red" {
		gameWinner = avalon.Red
	}

	if gameWinner == avalon.NoTeam {
		return fmt.Errorf("Game winner string could not be determined: %q\n", *winnerString)
	}

	basePath := BASE_PATH
	if *isTest {
		basePath = TEST_PATH
	}

	storage_, err := storage.LoadAvalonJsonStorageFromPath(basePath)
	if err != nil {
		fmt.Printf("avalon storage errored, making new storage")
		storage_, err = storage.NewAvalonJsonStorage(basePath)

		if err != nil {
			return errors.Wrapf(err, "NewAvalonJsonStorage(%q)", basePath)
		}
	}

	fmt.Printf("minions: %+v\n", minions)
	fmt.Printf("servants: %+v\n", loyalServants)

	roleByPlayerId := map[avalon.PlayerId]avalon.Role{}
	playerIds := []avalon.PlayerId{}

	// Add any new players if they don't exist
	for _, minionId := range *minions {
		_, err := storage_.GetPlayer(avalon.PlayerId(minionId))
		if err != nil {
			fmt.Printf("GetPlayer(%q) error: %s\n", minionId, err.Error())
		}

		if err != nil {
			if err = storage_.CreatePlayer(avalon.NewPlayer(minionId)); err != nil {
				return errors.Wrapf(err, "CreatePlayer(%q)", minionId)
			}
		}

		roleByPlayerId[avalon.PlayerId(minionId)] = avalon.MinionOfMordred
		playerIds = append(playerIds, avalon.PlayerId(minionId))
	}

	for _, servantId := range *loyalServants {
		_, err := storage_.GetPlayer(avalon.PlayerId(servantId))
		if err != nil {
			fmt.Printf("GetPlayer(%q) error: %s\n", servantId, err.Error())
		}

		if err != nil {
			if err = storage_.CreatePlayer(avalon.NewPlayer(servantId)); err != nil {
				return errors.Wrapf(err, "CreatePlayer(%q)", servantId)
			}
		}

		roleByPlayerId[avalon.PlayerId(servantId)] = avalon.LoyalServant
		playerIds = append(playerIds, avalon.PlayerId(servantId))
	}

	playersById, err := storage_.GetPlayersById(playerIds)
	if err != nil {
		return errors.Wrapf(err, "GetPlayersById")
	}

	game := avalon.NewGame(playersById, roleByPlayerId)
	game.SetWinner(gameWinner)

	updatedPlayers, err := game.UpdatePlayersAfterGame()
	if err != nil {
		return errors.Wrapf(err, "UpdatePlayersAfterGame")
	}

	if err = storage_.SaveGame(game); err != nil {
		return errors.Wrapf(err, "SaveGame")
	}

	if err = storage_.UpdatePlayers(updatedPlayers); err != nil {
		return errors.Wrapf(err, "UpdatePlayers")
	}

	return nil
}
