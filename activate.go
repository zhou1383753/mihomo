package main

import (
	"clash-admin/pkg/util"
	"net/http"
	"os"
)

func activate(w http.ResponseWriter, r *http.Request) {
	activationCode := r.FormValue("activationCode")
	if activationCode == "" {
		http.Error(w, "Invalid activation code", http.StatusBadRequest)
		return
	}

	if util.Decrypt([]byte(r.FormValue("macid")), activationCode) {
		err := os.WriteFile(".activation", []byte(activationCode), 0644)
		if err != nil {
			http.Error(w, "Error activating", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		http.Error(w, "Invalid activation code", http.StatusBadRequest)
		return
	}
}
