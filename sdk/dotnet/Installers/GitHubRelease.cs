// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Pde.Installers
{
    [PdeResourceType("pde:installers:GitHubRelease")]
    public partial class GitHubRelease : global::Pulumi.CustomResource
    {
        [Output("assetName")]
        public Output<string?> AssetName { get; private set; } = null!;

        [Output("download_url")]
        public Output<string> Download_url { get; private set; } = null!;

        [Output("environment")]
        public Output<ImmutableDictionary<string, string>?> Environment { get; private set; } = null!;

        [Output("executable")]
        public Output<string?> Executable { get; private set; } = null!;

        [Output("installCommands")]
        public Output<ImmutableArray<string>> InstallCommands { get; private set; } = null!;

        [Output("interpreter")]
        public Output<ImmutableArray<string>> Interpreter { get; private set; } = null!;

        [Output("org")]
        public Output<string> Org { get; private set; } = null!;

        [Output("releaseVersion")]
        public Output<string?> ReleaseVersion { get; private set; } = null!;

        [Output("repo")]
        public Output<string> Repo { get; private set; } = null!;

        [Output("uninstallCommands")]
        public Output<ImmutableArray<string>> UninstallCommands { get; private set; } = null!;

        [Output("updateCommands")]
        public Output<ImmutableArray<string>> UpdateCommands { get; private set; } = null!;

        [Output("version")]
        public Output<string> Version { get; private set; } = null!;


        /// <summary>
        /// Create a GitHubRelease resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public GitHubRelease(string name, GitHubReleaseArgs args, CustomResourceOptions? options = null)
            : base("pde:installers:GitHubRelease", name, args ?? new GitHubReleaseArgs(), MakeResourceOptions(options, ""))
        {
        }

        private GitHubRelease(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("pde:installers:GitHubRelease", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing GitHubRelease resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static GitHubRelease Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new GitHubRelease(name, id, options);
        }
    }

    public sealed class GitHubReleaseArgs : global::Pulumi.ResourceArgs
    {
        [Input("assetName")]
        public Input<string>? AssetName { get; set; }

        [Input("executable")]
        public Input<string>? Executable { get; set; }

        [Input("installCommands")]
        private InputList<string>? _installCommands;
        public InputList<string> InstallCommands
        {
            get => _installCommands ?? (_installCommands = new InputList<string>());
            set => _installCommands = value;
        }

        [Input("org", required: true)]
        public Input<string> Org { get; set; } = null!;

        [Input("releaseVersion")]
        public Input<string>? ReleaseVersion { get; set; }

        [Input("repo", required: true)]
        public Input<string> Repo { get; set; } = null!;

        [Input("uninstallCommands")]
        private InputList<string>? _uninstallCommands;
        public InputList<string> UninstallCommands
        {
            get => _uninstallCommands ?? (_uninstallCommands = new InputList<string>());
            set => _uninstallCommands = value;
        }

        [Input("updateCommands")]
        private InputList<string>? _updateCommands;
        public InputList<string> UpdateCommands
        {
            get => _updateCommands ?? (_updateCommands = new InputList<string>());
            set => _updateCommands = value;
        }

        public GitHubReleaseArgs()
        {
        }
        public static new GitHubReleaseArgs Empty => new GitHubReleaseArgs();
    }
}
