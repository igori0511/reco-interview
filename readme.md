// token name
// data-extractor

	// asana access token
	// 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3

	/*

				         curl https://app.asana.com/api/1.0/users/me \
								  -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"


							# API Endpoint: GET /users/me
							curl "https://app.asana.com/api/1.0/users/me" \
							     -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"


							all projects

						curl "https://app.asana.com/api/1.0/workspaces/1211834271836941/projects" \
						     -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"



					curl "https://app.asana.com/api/1.0/workspaces/1211834271836941/users" \
					     -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"


			# Replace 111111111 with your project_gid
			curl "https://app.asana.com/api/1.0/projects/1211834271836941/project_memberships" \
			     -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"


		//

		# This command gets all members for ONE specific project
		curl "https://app.asana.com/api/1.0/memberships?parent=1211834226422815" \
		     -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"
	*/