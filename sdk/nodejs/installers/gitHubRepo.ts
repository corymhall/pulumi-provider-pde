// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

export class GitHubRepo extends pulumi.CustomResource {
    /**
     * Get an existing GitHubRepo resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): GitHubRepo {
        return new GitHubRepo(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'pde:installers:GitHubRepo';

    /**
     * Returns true if the given object is an instance of GitHubRepo.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is GitHubRepo {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === GitHubRepo.__pulumiType;
    }

    public /*out*/ readonly absFolderName!: pulumi.Output<string>;
    public readonly branch!: pulumi.Output<string | undefined>;
    public /*out*/ readonly environment!: pulumi.Output<{[key: string]: string} | undefined>;
    public readonly folderName!: pulumi.Output<string | undefined>;
    public readonly installCommands!: pulumi.Output<string[] | undefined>;
    public /*out*/ readonly interpreter!: pulumi.Output<string[] | undefined>;
    public readonly org!: pulumi.Output<string>;
    public readonly repo!: pulumi.Output<string>;
    public readonly uninstallCommands!: pulumi.Output<string[] | undefined>;
    public readonly updateCommands!: pulumi.Output<string[] | undefined>;
    public /*out*/ readonly version!: pulumi.Output<string>;

    /**
     * Create a GitHubRepo resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: GitHubRepoArgs, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.org === undefined) && !opts.urn) {
                throw new Error("Missing required property 'org'");
            }
            if ((!args || args.repo === undefined) && !opts.urn) {
                throw new Error("Missing required property 'repo'");
            }
            resourceInputs["branch"] = args ? args.branch : undefined;
            resourceInputs["folderName"] = args ? args.folderName : undefined;
            resourceInputs["installCommands"] = args ? args.installCommands : undefined;
            resourceInputs["org"] = args ? args.org : undefined;
            resourceInputs["repo"] = args ? args.repo : undefined;
            resourceInputs["uninstallCommands"] = args ? args.uninstallCommands : undefined;
            resourceInputs["updateCommands"] = args ? args.updateCommands : undefined;
            resourceInputs["absFolderName"] = undefined /*out*/;
            resourceInputs["environment"] = undefined /*out*/;
            resourceInputs["interpreter"] = undefined /*out*/;
            resourceInputs["version"] = undefined /*out*/;
        } else {
            resourceInputs["absFolderName"] = undefined /*out*/;
            resourceInputs["branch"] = undefined /*out*/;
            resourceInputs["environment"] = undefined /*out*/;
            resourceInputs["folderName"] = undefined /*out*/;
            resourceInputs["installCommands"] = undefined /*out*/;
            resourceInputs["interpreter"] = undefined /*out*/;
            resourceInputs["org"] = undefined /*out*/;
            resourceInputs["repo"] = undefined /*out*/;
            resourceInputs["uninstallCommands"] = undefined /*out*/;
            resourceInputs["updateCommands"] = undefined /*out*/;
            resourceInputs["version"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(GitHubRepo.__pulumiType, name, resourceInputs, opts);
    }
}

/**
 * The set of arguments for constructing a GitHubRepo resource.
 */
export interface GitHubRepoArgs {
    branch?: pulumi.Input<string>;
    folderName?: pulumi.Input<string>;
    installCommands?: pulumi.Input<pulumi.Input<string>[]>;
    org: pulumi.Input<string>;
    repo: pulumi.Input<string>;
    uninstallCommands?: pulumi.Input<pulumi.Input<string>[]>;
    updateCommands?: pulumi.Input<pulumi.Input<string>[]>;
}
