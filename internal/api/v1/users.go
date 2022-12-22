package v1

import (
	"encoding/json"
	"net/http"

	"homework/internal/util/terr"
	"homework/specs"
)

func (a apiServer) GetUserById(w http.ResponseWriter, r *http.Request, userIdSpecs specs.UUIDPathObjectID) {

	userId, err := convertStringToUuid(string(userIdSpecs))
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_USER_UUID", err.Error()))
		return
	}

	ctx := r.Context()
	user, err := a.serviceRegistry.User.GetUserById(ctx, userId)

	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	userSpecs := transformUser(user)
	_ = json.NewEncoder(w).Encode(userSpecs)

}
