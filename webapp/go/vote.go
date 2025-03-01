package main

import "strings"

// Vote Model
type Vote struct {
	ID          int
	UserID      int
	CandidateID int
	Keyword     string
}

func getVoteCountByCandidateID(candidateID int) (count int) {
	row := db.QueryRow("SELECT COUNT(id) AS count FROM votes WHERE candidate_id = ?", candidateID)
	row.Scan(&count)
	return
}

func getUserVotedCount(userID int) (count int) {
	row := db.QueryRow("SELECT COUNT(id) AS count FROM votes WHERE user_id = ?", userID)
	row.Scan(&count)
	return
}

func createVote(userID int, candidateID int, keyword string) {
	db.Exec("INSERT INTO votes (user_id, candidate_id, keyword) VALUES (?, ?, ?)",
		userID, candidateID, keyword)
}

func getVoiceOfSupporter(candidateIDs []int) (voices []string) {
	args := []interface{}{}
	for _, candidateID := range candidateIDs {
		args = append(args, candidateID)
	}
	rows, err := db.Query(`
    SELECT keyword
    FROM votes
    WHERE candidate_id IN (`+strings.Join(strings.Split(strings.Repeat("?", len(candidateIDs)), ""), ",")+`)
    GROUP BY keyword
    ORDER BY COUNT(*) DESC
    LIMIT 10`, args...)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		var keyword string
		err = rows.Scan(&keyword)
		if err != nil {
			panic(err.Error())
		}
		voices = append(voices, keyword)
	}
	return
}
