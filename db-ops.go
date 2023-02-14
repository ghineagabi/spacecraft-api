package main

import (
	"errors"
	"fmt"
	"strings"
)

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
