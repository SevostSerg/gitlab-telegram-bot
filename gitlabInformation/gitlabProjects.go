package gitlabInformation

import (
	config "GitlabTgBot/configuration"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GitlabProject struct {
	ID                                        int64                     `json:"id"`
	Description                               DescriptionEnum           `json:"description"`
	Name                                      string                    `json:"name"`
	NameWithNamespace                         string                    `json:"name_with_namespace"`
	Path                                      string                    `json:"path"`
	PathWithNamespace                         string                    `json:"path_with_namespace"`
	CreatedAt                                 string                    `json:"created_at"`
	TagList                                   []interface{}             `json:"tag_list"`
	Topics                                    []interface{}             `json:"topics"`
	SSHURLToRepo                              string                    `json:"ssh_url_to_repo"`
	HTTPURLToRepo                             string                    `json:"http_url_to_repo"`
	WebURL                                    string                    `json:"web_url"`
	ReadmeURL                                 *string                   `json:"readme_url"`
	AvatarURL                                 interface{}               `json:"avatar_url"`
	ForksCount                                int64                     `json:"forks_count"`
	StarCount                                 int64                     `json:"star_count"`
	LastActivityAt                            string                    `json:"last_activity_at"`
	Namespace                                 Namespace                 `json:"namespace"`
	ContainerRegistryImagePrefix              string                    `json:"container_registry_image_prefix"`
	Links                                     Links                     `json:"_links"`
	PackagesEnabled                           bool                      `json:"packages_enabled"`
	EmptyRepo                                 bool                      `json:"empty_repo"`
	Archived                                  bool                      `json:"archived"`
	Visibility                                PagesAccessLevel          `json:"visibility"`
	ResolveOutdatedDiffDiscussions            bool                      `json:"resolve_outdated_diff_discussions"`
	ContainerExpirationPolicy                 ContainerExpirationPolicy `json:"container_expiration_policy"`
	IssuesEnabled                             bool                      `json:"issues_enabled"`
	MergeRequestsEnabled                      bool                      `json:"merge_requests_enabled"`
	WikiEnabled                               bool                      `json:"wiki_enabled"`
	JobsEnabled                               bool                      `json:"jobs_enabled"`
	SnippetsEnabled                           bool                      `json:"snippets_enabled"`
	ContainerRegistryEnabled                  bool                      `json:"container_registry_enabled"`
	ServiceDeskEnabled                        bool                      `json:"service_desk_enabled"`
	ServiceDeskAddress                        interface{}               `json:"service_desk_address"`
	CanCreateMergeRequestIn                   bool                      `json:"can_create_merge_request_in"`
	IssuesAccessLevel                         AnalyticsAccessLevel      `json:"issues_access_level"`
	RepositoryAccessLevel                     AnalyticsAccessLevel      `json:"repository_access_level"`
	MergeRequestsAccessLevel                  AnalyticsAccessLevel      `json:"merge_requests_access_level"`
	ForkingAccessLevel                        AnalyticsAccessLevel      `json:"forking_access_level"`
	WikiAccessLevel                           AnalyticsAccessLevel      `json:"wiki_access_level"`
	BuildsAccessLevel                         AnalyticsAccessLevel      `json:"builds_access_level"`
	SnippetsAccessLevel                       AnalyticsAccessLevel      `json:"snippets_access_level"`
	PagesAccessLevel                          PagesAccessLevel          `json:"pages_access_level"`
	OperationsAccessLevel                     AnalyticsAccessLevel      `json:"operations_access_level"`
	AnalyticsAccessLevel                      AnalyticsAccessLevel      `json:"analytics_access_level"`
	ContainerRegistryAccessLevel              AnalyticsAccessLevel      `json:"container_registry_access_level"`
	EmailsDisabled                            *bool                     `json:"emails_disabled"`
	SharedRunnersEnabled                      bool                      `json:"shared_runners_enabled"`
	LFSEnabled                                bool                      `json:"lfs_enabled"`
	CreatorID                                 int64                     `json:"creator_id"`
	ImportStatus                              ImportStatus              `json:"import_status"`
	CiDefaultGitDepth                         int64                     `json:"ci_default_git_depth"`
	CiForwardDeploymentEnabled                bool                      `json:"ci_forward_deployment_enabled"`
	CiJobTokenScopeEnabled                    bool                      `json:"ci_job_token_scope_enabled"`
	PublicJobs                                bool                      `json:"public_jobs"`
	BuildTimeout                              int64                     `json:"build_timeout"`
	AutoCancelPendingPipelines                AnalyticsAccessLevel      `json:"auto_cancel_pending_pipelines"`
	BuildCoverageRegex                        interface{}               `json:"build_coverage_regex"`
	SharedWithGroups                          []SharedWithGroup         `json:"shared_with_groups"`
	OnlyAllowMergeIfPipelineSucceeds          bool                      `json:"only_allow_merge_if_pipeline_succeeds"`
	AllowMergeOnSkippedPipeline               interface{}               `json:"allow_merge_on_skipped_pipeline"`
	RestrictUserDefinedVariables              bool                      `json:"restrict_user_defined_variables"`
	RequestAccessEnabled                      bool                      `json:"request_access_enabled"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool                      `json:"only_allow_merge_if_all_discussions_are_resolved"`
	RemoveSourceBranchAfterMerge              bool                      `json:"remove_source_branch_after_merge"`
	PrintingMergeRequestLinkEnabled           bool                      `json:"printing_merge_request_link_enabled"`
	MergeMethod                               MergeMethod               `json:"merge_method"`
	SquashOption                              SquashOption              `json:"squash_option"`
	SuggestionCommitMessage                   interface{}               `json:"suggestion_commit_message"`
	AutoDevopsEnabled                         bool                      `json:"auto_devops_enabled"`
	AutoDevopsDeployStrategy                  AutoDevopsDeployStrategy  `json:"auto_devops_deploy_strategy"`
	AutocloseReferencedIssues                 bool                      `json:"autoclose_referenced_issues"`
	KeepLatestArtifact                        bool                      `json:"keep_latest_artifact"`
	Permissions                               Permissions               `json:"permissions"`
	DefaultBranch                             *string                   `json:"default_branch,omitempty"`
	Owner                                     *GitlabUser               `json:"owner,omitempty"`
	CiConfigPath                              *string                   `json:"ci_config_path,omitempty"`
}

type ContainerExpirationPolicy struct {
	Cadence       Cadence     `json:"cadence"`
	Enabled       bool        `json:"enabled"`
	KeepN         int64       `json:"keep_n"`
	OlderThan     OlderThan   `json:"older_than"`
	NameRegex     NameRegex   `json:"name_regex"`
	NameRegexKeep interface{} `json:"name_regex_keep"`
	NextRunAt     string      `json:"next_run_at"`
}

type Links struct {
	Self          string `json:"self"`
	MergeRequests string `json:"merge_requests"`
	RepoBranches  string `json:"repo_branches"`
	Labels        string `json:"labels"`
	Events        string `json:"events"`
	Members       string `json:"members"`
}

type Namespace struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Kind      Kind        `json:"kind"`
	FullPath  string      `json:"full_path"`
	ParentID  interface{} `json:"parent_id"`
	AvatarURL *string     `json:"avatar_url"`
	WebURL    string      `json:"web_url"`
}

type Permissions struct {
	ProjectAccess interface{}  `json:"project_access"`
	GroupAccess   *GroupAccess `json:"group_access"`
}

type GroupAccess struct {
	AccessLevel       int64 `json:"access_level"`
	NotificationLevel int64 `json:"notification_level"`
}

type SharedWithGroup struct {
	GroupID          int64       `json:"group_id"`
	GroupName        string      `json:"group_name"`
	GroupFullPath    string      `json:"group_full_path"`
	GroupAccessLevel int64       `json:"group_access_level"`
	ExpiresAt        interface{} `json:"expires_at"`
}

type AnalyticsAccessLevel string

const (
	Disabled AnalyticsAccessLevel = "disabled"
	Enabled  AnalyticsAccessLevel = "enabled"
)

type AutoDevopsDeployStrategy string

const (
	Continuous AutoDevopsDeployStrategy = "continuous"
)

type Cadence string

const (
	The1D Cadence = "1d"
)

type NameRegex string

const (
	Empty NameRegex = ".*"
)

type OlderThan string

const (
	The90D OlderThan = "90d"
)

type DescriptionEnum string

const (
	ContainsASimpleResponderClientDoNotUseInProduction DescriptionEnum = "Contains a simple responder client. Do not use in production."
	ContainsASimpleResponderDoNotUseInProduction       DescriptionEnum = "Contains a simple responder. Do not use in production."
	Description                                        DescriptionEnum = ""
)

type ImportStatus string

const (
	None ImportStatus = "none"
)

type MergeMethod string

const (
	Merge MergeMethod = "merge"
)

type Kind string

const (
	Group Kind = "group"
	User  Kind = "user"
)

type PagesAccessLevel string

const (
	Internal PagesAccessLevel = "internal"
	Private  PagesAccessLevel = "private"
)

type SquashOption string

const (
	DefaultOff SquashOption = "default_off"
)

// projects limit = 100
const (
	requestAddition string = "/api/v4/projects?per_page=100&private_token="
)

func GetProjectsList(accessToken string) ([]GitlabProject, error) {
	URL := config.GetConfigInstance().GitlabURL
	response, err := http.Get(URL + requestAddition + accessToken)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
	}

	var projectList []GitlabProject
	err = json.Unmarshal(contents, &projectList)
	if err != nil {
		log.Panic(err)
	}

	return projectList, nil
}
