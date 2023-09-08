// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

export class GitHubRelease extends pulumi.CustomResource {
    /**
     * Get an existing GitHubRelease resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): GitHubRelease {
        return new GitHubRelease(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'pde:installers:GitHubRelease';

    /**
     * Returns true if the given object is an instance of GitHubRelease.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is GitHubRelease {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === GitHubRelease.__pulumiType;
    }

    public readonly assetName!: pulumi.Output<string | undefined>;
    public readonly binFolder!: pulumi.Output<string | undefined>;
    public readonly binLocation!: pulumi.Output<string | undefined>;
    public /*out*/ readonly downloadURL!: pulumi.Output<string>;
    public /*out*/ readonly environment!: pulumi.Output<{[key: string]: string} | undefined>;
    public readonly executable!: pulumi.Output<string | undefined>;
    public readonly installCommands!: pulumi.Output<string[] | undefined>;
    public /*out*/ readonly interpreter!: pulumi.Output<string[] | undefined>;
    public /*out*/ readonly locations!: pulumi.Output<string[] | undefined>;
    public readonly org!: pulumi.Output<string>;
    public readonly repo!: pulumi.Output<string>;
    public readonly uninstallCommands!: pulumi.Output<string[] | undefined>;
    public readonly updateCommands!: pulumi.Output<string[] | undefined>;
    public readonly version!: pulumi.Output<string | undefined>;

    /**
     * Create a GitHubRelease resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: GitHubReleaseArgs, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.org === undefined) && !opts.urn) {
                throw new Error("Missing required property 'org'");
            }
            if ((!args || args.repo === undefined) && !opts.urn) {
                throw new Error("Missing required property 'repo'");
            }
            resourceInputs["assetName"] = args ? args.assetName : undefined;
            resourceInputs["binFolder"] = args ? args.binFolder : undefined;
            resourceInputs["binLocation"] = args ? args.binLocation : undefined;
            resourceInputs["executable"] = args ? args.executable : undefined;
            resourceInputs["installCommands"] = args ? args.installCommands : undefined;
            resourceInputs["org"] = args ? args.org : undefined;
            resourceInputs["repo"] = args ? args.repo : undefined;
            resourceInputs["uninstallCommands"] = args ? args.uninstallCommands : undefined;
            resourceInputs["updateCommands"] = args ? args.updateCommands : undefined;
            resourceInputs["version"] = args ? args.version : undefined;
            resourceInputs["downloadURL"] = undefined /*out*/;
            resourceInputs["environment"] = undefined /*out*/;
            resourceInputs["interpreter"] = undefined /*out*/;
            resourceInputs["locations"] = undefined /*out*/;
        } else {
            resourceInputs["assetName"] = undefined /*out*/;
            resourceInputs["binFolder"] = undefined /*out*/;
            resourceInputs["binLocation"] = undefined /*out*/;
            resourceInputs["downloadURL"] = undefined /*out*/;
            resourceInputs["environment"] = undefined /*out*/;
            resourceInputs["executable"] = undefined /*out*/;
            resourceInputs["installCommands"] = undefined /*out*/;
            resourceInputs["interpreter"] = undefined /*out*/;
            resourceInputs["locations"] = undefined /*out*/;
            resourceInputs["org"] = undefined /*out*/;
            resourceInputs["repo"] = undefined /*out*/;
            resourceInputs["uninstallCommands"] = undefined /*out*/;
            resourceInputs["updateCommands"] = undefined /*out*/;
            resourceInputs["version"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(GitHubRelease.__pulumiType, name, resourceInputs, opts);
    }
}

/**
 * The set of arguments for constructing a GitHubRelease resource.
 */
export interface GitHubReleaseArgs {
    assetName?: pulumi.Input<string>;
    binFolder?: pulumi.Input<string>;
    binLocation?: pulumi.Input<string>;
    executable?: pulumi.Input<string>;
    installCommands?: pulumi.Input<pulumi.Input<string>[]>;
    org: pulumi.Input<string>;
    repo: pulumi.Input<string>;
    uninstallCommands?: pulumi.Input<pulumi.Input<string>[]>;
    updateCommands?: pulumi.Input<pulumi.Input<string>[]>;
    version?: pulumi.Input<string>;
}
