# Backend Bucket - Terraform

## Setup

<walkthrough-author name="rileykarson@google.com" analyticsId="UA-125550242-1" tutorialName="bbdev" repositoryUrl="https://github.com/rileykarson/terraform-provider-google"></walkthrough-author>

Welcome to Terraform in Cloud Shell! We need you to let us know what project you'd like to use with Terraform.

<walkthrough-project-billing-setup permissions="compute.backendBuckets.create"></walkthrough-project-billing-setup>

Terraform uses real resources, so anything you create in this session will be billed against this project.

This tutorial requires the Compute API. Not sure if you have it enabled? Click the "Enable APIs" button below.
Otherwise, go ahead and click "Forward"!

<walkthrough-enable-apis apis="compute.googleapis.com"></walkthrough-enable-apis>

## Terraforming!

Let's get Terraform using {{project-id}}! Click the Cloud Shell icon below to copy the command
to your shell, and then run it from the shell by pressing Enter/Return. Terraform will pick up
the project name from the environment variable.

```bash
export GOOGLE_CLOUD_PROJECT={{project-id}}
```

After that, let's get Terraform started. Run the following to pull in the providers.

```bash
terraform init
```

With the providers and a project, you're ready to use Terraform. Go ahead!

```bash
terraform apply
```

Terraform will show you what it plans to do, and prompt you to accept. Type "yes" to accept the plan.

```bash
yes
```


## Post-Apply

### Editing your config

Now you've provisioned your resources in GCP! If you run a "plan", you should see no changes needed.

```bash
terraform plan
```

So let's make a change! Try editing a number, or appending a value to the name in the editor. Then,
run a 'plan' again.

```bash
terraform plan
```

Afterwards you can run an apply, which implicitly does a plan and shows you the intended changes
at the 'yes' prompt.

```bash
terraform apply
```

```bash
yes
```

### Cleanup

Run the following to remove the resources Terraform provisioned:

```bash
terraform destroy
```
```bash
yes
```
