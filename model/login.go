package model

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
    Token  string `json:"token"`
    ClientManagerAssociations []model.ClientManagerAssociation `json:"client_manager_associations"`
}
