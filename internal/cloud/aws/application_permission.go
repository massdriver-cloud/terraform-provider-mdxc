package aws

// type IAMAPI interface {
// 	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
// 	DeleteRole(ctx context.Context, params *iam.DeleteRoleInput, optFns ...func(*iam.Options)) (*iam.DeleteRoleOutput, error)
// }

// // TODO: call this and inject into Create()
// func (c AWSConfig) NewIAMService() IAMAPI {
// 	client := iam.NewFromConfig(*c.config)
// 	return client
// }

// // Create an AWS IAM Role as a massdriver.AppIdentity
// func AppIdentityCreate(ctx context.Context, api IAMCreateRoleAPI, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
// 	roleInput := iam.CreateRoleInput{
// 		RoleName: input.Name,
// 	}

// 	roleOutput, err := api.CreateRole(ctx, &roleInput, []func(*iam.Options){})

// 	appIdentityOutput := massdriver.AppIdentityOutput{
// 		AwsIamRole: *roleOutput.Role,
// 	}

// 	return &appIdentityOutput, err
// }

// type applicationIdentityConfig struct {
// 	Name             string
// 	AssumeRolePolicy string
// }

// func (c *AWSConfig) CreateApplicationPermission(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
// 	iamClient := c.NewIAMService()

// 	applicationIdentityConfig, extractErr := extractApplicationIdentityConfig(d)
// 	if extractErr != nil {
// 		diag.FromErr(extractErr)
// 	}

// 	output, roleErr := doCreateAWSIAMRole(ctx, applicationIdentityConfig, iamClient)
// 	if roleErr != nil {
// 		return diag.FromErr(roleErr)
// 	}

// 	roleName := *output.Role.RoleName
// 	d.SetId(roleName)

// 	return nil
// }

// func (c *AWSConfig) DeleteApplicationPermission(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
// 	iamClient := c.NewIAMService()

// 	applicationIdentityConfig, extractErr := extractApplicationIdentityConfig(d)
// 	if extractErr != nil {
// 		diag.FromErr(extractErr)
// 	}

// 	_, roleErr := doDeleteAWSIAMRole(ctx, applicationIdentityConfig, iamClient)
// 	if roleErr != nil {
// 		return diag.FromErr(roleErr)
// 	}

// 	return nil
// }

// func extractApplicationIdentityConfig(d *schema.ResourceData) (*applicationIdentityConfig, error) {
// 	applicationIdentityConfig := applicationIdentityConfig{}

// 	if v, ok := d.Get("name").(string); ok && v != "" {
// 		applicationIdentityConfig.Name = v
// 	}

// 	var awsNestedMap map[string]interface{}
// 	if awsBlock, ok := d.Get("aws").([]interface{}); ok && len(awsBlock) > 0 && awsBlock[0] != nil {
// 		awsNestedMap = awsBlock[0].(map[string]interface{})
// 	} else {
// 		return nil, errors.New("AWS configuration block not specified")
// 	}

// 	if v, ok := awsNestedMap["assume_role_policy"].(string); ok && v != "" {
// 		applicationIdentityConfig.AssumeRolePolicy = v
// 	}

// 	return &applicationIdentityConfig, nil
// }

// func doCreateAWSIAMRole(ctx context.Context, config *applicationIdentityConfig, client IAMAPI) (*iam.CreateRoleOutput, error) {

// 	assumeRolePolicy, assumeErr := structure.NormalizeJsonString(config.AssumeRolePolicy)
// 	if assumeErr != nil {
// 		return nil, assumeErr
// 	}

// 	roleInput := iam.CreateRoleInput{
// 		AssumeRolePolicyDocument: &assumeRolePolicy,
// 		RoleName:                 &config.Name,
// 	}

// 	return client.CreateRole(ctx, &roleInput)
// }

// func doDeleteAWSIAMRole(ctx context.Context, config *applicationIdentityConfig, client IAMAPI) (*iam.DeleteRoleOutput, error) {

// 	input := iam.DeleteRoleInput{
// 		RoleName: aws.String(config.Name),
// 	}

// 	return client.DeleteRole(ctx, &input)
// }
