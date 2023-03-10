package main

import (
	"errors"
	"fmt"
	"strings"
)

// CheckCredentials queries the DB with the given user credentials and returns an error if 0 results are found
// Note: The password in the DB should be saved as a SHA512
// Credentials (in Basic Auth): root: admin
// Credentials (in DB): root: 30bb8411dd0cbf96b10a52371f7b3be1690f7afa16c3bd7bc7d02c0e2854768d
func CheckCredentials(u *UserCredentials) error {
	sqlStatement := `SELECT * FROM new_schema.user WHERE email = ? AND password = ?`
	row, err := Db.Query(sqlStatement, u.Email, SHA512(u.Password))
	defer row.Close()
	if err != nil {
		return err
	}
	if !row.Next() {
		return &InvalidFieldsError{Location: "Basic auth", AffectedField: "email and password", Reason: "EmailID and/or password mismatch"}
	}

	if err != nil {
		return err
	}
	return nil
}

func InsertSpacecraft(sp *CreateSpacecraft) error {
	sqlStatement := `INSERT INTO new_schema.spacecraft (name, class, crew, image, value, status) 
VALUES (?, ?, ?, ?, ?, ?);`

	_, err := Db.Exec(sqlStatement, sp.Name, sp.Class,
		sp.Crew, sp.Image, sp.Value, sp.Status)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSpacecraftFromDB(emailID int) error {
	sqlStatement := `DELETE FROM new_schema.spacecraft WHERE id = ?;`
	res, err := Db.Exec(sqlStatement, emailID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

// UpdateSpacecraftFromDB updates the DB fields with the provided UpdateSpacecraft parameters
func UpdateSpacecraftFromDB(u *UpdateSpacecraft) error {
	var query strings.Builder
	params := make([]interface{}, 0)
	params = append(params)
	query.WriteString("UPDATE new_schema.spacecraft SET")

	if u.Name != "" {
		query.WriteString(fmt.Sprintf(" name=?,"))
		params = append(params, u.Name)
	}
	if u.Image != "" {
		query.WriteString(fmt.Sprintf(" image=?,"))
		params = append(params, u.Image)
	}
	if u.Class != "" {
		query.WriteString(fmt.Sprintf(" class=?,"))
		params = append(params, u.Class)
	}
	if u.Crew != 0 {
		query.WriteString(fmt.Sprintf(" crew=?,"))
		params = append(params, u.Crew)
	}
	if u.Value != 0 {
		query.WriteString(fmt.Sprintf(" value=?,"))
		params = append(params, u.Value)
	}
	if u.Status != "" {
		query.WriteString(fmt.Sprintf(" status=?,"))
		params = append(params, u.Status)
	}

	if len(params) < 1 {
		return errors.New(NOTENOUGHPARAMS)
	}

	queryString := fmt.Sprintf("%s WHERE ID = %d", strings.TrimSuffix(query.String(), ","), u.Id)

	_, err := Db.Exec(queryString, params...)
	if err != nil {
		return err
	}
	return nil
}

/*
	GetSpacecraftsFromDB makes 2 queries:

1) Selects the armaments belonging to the provided ID
2) Selects the spacecraft having the provided ID
Then Iterates through the queried armaments and appends them to the DetailedSpacecraft struct
*/
func GetSpacecraftsFromDB(s *DetailedSpacecraft, id int) error {
	sqlStatementForArmaments := `SELECT title, qty FROM new_schema.armaments
WHERE armaments.id_spacecraft = ?`

	rows, err := Db.Query(sqlStatementForArmaments, id)
	defer rows.Close()

	if err != nil {
		return err
	}

	var armaments []Armament
	var arm Armament
	for rows.Next() {
		if err = rows.Scan(&arm.Title, &arm.Qty); err != nil {
			return err
		}
		armaments = append(armaments, arm)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	sqlStatementForSpacecraft := `SELECT id, name, class, crew, image, value, status
FROM new_schema.spacecraft
WHERE spacecraft.id = ?`

	rows, err = Db.Query(sqlStatementForSpacecraft, id)
	defer rows.Close()

	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.Scan(&s.Id, &s.Name, &s.Class, &s.Crew, &s.Image, &s.Value, &s.Status); err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	for _, a := range armaments {
		s.Armament = append(s.Armament, a)
	}

	return nil
}

// GetFilteredSpacecrafts queries the DB based on the <name, class, status> parameters.
// It finds if the provided <name, class, status> substrings are inside the <name, class, status> fields
/* Ex :(DB)
{Spacecraft: {name: Torch, class: High, status: Up},
 Spacecraft: {name: Born,  class: Hi, status: Dump}}

   (Query) : {name: "or", class: "i", status: "p"}

Would return both the spacecrafts because "or" is in both names, "i" is in both classes and "p" is in both statuses
*/
func GetFilteredSpacecrafts(s *[]FilteredSpacecraft, name, class, status string) error {
	sqlStatement := `SELECT id, name, status
FROM new_schema.spacecraft
WHERE LOWER(name) LIKE concat('%', ?, '%') AND 
      LOWER(class) LIKE concat('%', ?, '%') AND 
      LOWER(status) LIKE concat('%', ?, '%')`
	rows, err := Db.Query(sqlStatement, name, class, status)
	defer rows.Close()

	var spacecraft FilteredSpacecraft
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.Scan(&spacecraft.Id, &spacecraft.Name, &spacecraft.Status); err != nil {
			return err
		}
		*s = append(*s, spacecraft)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
