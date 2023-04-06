package main

import (
	"database/sql"
	"fmt"
	"github.com/rmarken5/lava/inspect/data-access"
	"github.com/rmarken5/lava/inspect/logic"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	// Define the connection parameters
	connStr := "user=user password=password dbname=mlb host=localhost port=5432 sslmode=disable"

	// Create the connection pool
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
	}
	defer db.Close()

	inspector := data_access.NewInspector(db, log.New(NewStdoutWriter(), "data-access: ", 0))

	logic := logic.New(log.New(NewStdoutWriter(), "logic: ", 0), inspector)

	str, err := logic.BuildStructsForQuery(`select g.*,
       gds.id as "defense.id",
       gds.team_id as "defense.team_id",
       gds.game_id as "defense.game_id",
       gds.defense_assists as "defense.defense_assists",
       gds.defense_caught_stealing as "defense.defense_caught_stealing",
       gds.defense_chances as "defense.defense_chances",
       gds.defense_errors as "defense.defense_errors",
       gds.defense_passed_balls as "defense.defense_passed_balls",
       gds.defense_pick_offs as "defense.defense_pick_offs",
       gds.defense_put_outs as "defense.defense_put_outs",
       gds.defense_runs as "defense.defense_runs",
       gds.created_at as "defense.created_at",
       gds.updated_at as "defense.updated_at",
       gds.deleted_at as "defense.deleted_at",
       gos.id as  "offense.id",
       gos.team_id as "offense.team_id",
       gos.game_id as "offense.game_id",
       gos.offense_at_bats as "offense.offense_at_bats",
       gos.offense_avg as "offense.offense_avg",
       gos.offense_base_on_balls as "offense.offense_base_on_balls",
       gos.offense_caught_stealing as "offense.offense_caught_stealing",
       gos.offense_doubles as "offense.offense_doubles",
       gos.offense_flyouts as "offense.offense_flyouts",
       gos.offense_ground_into_dp as "offense.offense_ground_into_dp",
       gos.offense_ground_into_tp as "offense.offense_ground_into_tp",
       gos.offense_ground_outs as "offense.offense_ground_outs",
       gos.offense_hbp as "offense.offense_hbp",
       gos.offense_hits as "offense.offense_hits",
       gos.offense_hr as "offense.offense_hr",
       gos.offense_intentional_walks as "offense.offense_intentional_walks",
       gos.offense_left_on_base as "offense.offense_left_on_base",
       gos.offense_obp as "offense.offense_obp",
       gos.offense_ops as "offense.offense_ops",
       gos.offense_pick_offs as "offense.offense_pick_offs",
       gos.offense_plate_appearances as "offense.offense_plate_appearances",
       gos.offense_rbi as "offense.offense_rbi",
       gos.offense_runs as "offense.offense_runs",
       gos.offense_sac_bunts as "offense.offense_sac_bunts",
       gos.offense_sac_fly as "offense.offense_sac_fly",
       gos.offense_stolen_bases as "offense.offense_stolen_bases",
       gos.offense_strike_outs as "offense.offense_strike_outs",
       gos.offense_triples as "offense.offense_triples",
       gos.created_at as "offense.created_at",
       gos.updated_at as "offense.updated_at",
       gos.deleted_at as "offense.deleted_at",
       gps.id as "pitching.id",
       gps.team_id as "pitching.team_id",
       gps.game_id as "pitching.game_id",
       gps.pitching_air_outs as "pitching.pitching_air_outs",
       gps.pitching_at_bats as "pitching.pitching_at_bats",
       gps.pitching_balks as "pitching.pitching_balks",
       gps.pitching_balls as "pitching.pitching_balls",
       gps.pitching_base_on_balls as "pitching.pitching_base_on_balls",
       gps.pitching_batters_faced as "pitching.pitching_batters_faced",
       gps.pitching_caught_stealing as "pitching.pitching_caught_stealing",
       gps.pitching_complete_game as "pitching.pitching_complete_game",
       gps.pitching_doubles as "pitching.pitching_doubles",
       gps.pitching_earned_runs as "pitching.pitching_earned_runs",
       gps.pitching_earned_run_avg as "pitching.pitching_earned_run_avg",
       gps.pitching_ground_outs as "pitching.pitching_ground_outs",
       gps.pitching_hit_batsman as "pitching.pitching_hit_batsman",
       gps.pitching_hit_by_pitch as "pitching.pitching_hit_by_pitch",
       gps.pitching_hits as "pitching.pitching_hits",
       gps.pitching_home_runs as "pitching.pitching_home_runs",
       gps.pitching_innings_pitched as "pitching.pitching_innings_pitched",
       gps.pitching_intentional_walks as "pitching.pitching_intentional_walks",
       gps.pitching_num_of_pitches as "pitching.pitching_num_of_pitches",
       gps.pitching_obp as "pitching.pitching_obp",
       gps.pitching_outs as "pitching.pitching_outs",
       gps.pitching_passed_balls as "pitching.pitching_passed_balls",
       gps.pitching_rbi as "pitching.pitching_rbi",
       gps.pitching_runs as "pitching.pitching_runs",
       gps.pitching_sac_bunts as "pitching.pitching_sac_bunts",
       gps.pitching_stolen_bases as "pitching.pitching_stolen_bases",
       gps.pitching_strike_outs as "pitching.pitching_strike_outs",
       gps.pitching_strikes as "pitching.pitching_strikes",
       gps.pitching_triples as "pitching.pitching_triples",
       gps.pitching_whip as "pitching.pitching_whip",
       gps.pitching_wild_pitches as "pitching.pitching_wild_pitches",
       gps.created_at as "pitching.created_at",
       gps.updated_at as "pitching.updated_at",
       gps.deleted_at as "pitching.deleted_at"
from game as g
         inner join game_defensive_stats as gds on g.id = gds.game_id
         inner join game_offensive_stats as gos on gds.game_id = gos.game_id and gds.team_id = gos.team_id
         inner join game_pitching_stats as gps on gos.game_id = gps.game_id and gos.team_id = gps.team_id;`)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
	}
	fmt.Printf("%s", str)

}

func NewStdoutWriter() io.Writer {
	return os.Stdout
}
