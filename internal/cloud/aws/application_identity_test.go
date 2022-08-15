package aws_test

// type mockCreateRoleAPI func(ctx context.Context, params *iam.CreateRoleInput, optFns []func(*iam.Options)) (*iam.CreateRoleOutput, error)

// func (m mockCreateRoleAPI) CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns []func(*iam.Options)) (*iam.CreateRoleOutput, error) {
// 	return m(ctx, params, optFns)
// }

// func TestCreate(t *testing.T) {
// 	m := mockCreateRoleAPI(func(ctx context.Context, params *iam.CreateRoleInput, optFns []func(*iam.Options)) (*iam.CreateRoleOutput, error) {
// 		t.Helper()
// 		if params.RoleName == nil {
// 			t.Fatal("expect role name to not be nil")
// 		}

// 		arn := fmt.Sprintf("arn:aws:iam::account:role/%s", *params.RoleName)
// 		role := &types.Role{Arn: aws.String(arn)}
// 		return &iam.CreateRoleOutput{Role: role}, nil
// 	})

// 	appIdentityInput := massdriver.AppIdentityInput{
// 		Name: aws.String("test"),
// 	}

// 	appIdentityOutput, _ := mdxcaws.AppIdentityCreate(context.TODO(), m, &appIdentityInput)
// 	got := *appIdentityOutput.AwsIamRole.Arn
// 	want := "arn:aws:iam::account:role/test"

// 	if want != got {
// 		t.Errorf("expect %v, got %v", want, got)
// 	}
// }

// TODO: Add TestCreateError case
