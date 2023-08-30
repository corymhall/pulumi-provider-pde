// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

export class Link extends pulumi.CustomResource {
    /**
     * Get an existing Link resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): Link {
        return new Link(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'pde:local:Link';

    /**
     * Returns true if the given object is an instance of Link.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Link {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Link.__pulumiType;
    }

    public readonly exists!: pulumi.Output<boolean>;
    public readonly is_dir!: pulumi.Output<boolean>;
    public readonly linked!: pulumi.Output<boolean>;
    public readonly overwrite!: pulumi.Output<boolean>;
    public /*out*/ readonly result!: pulumi.Output<string>;
    public readonly source!: pulumi.Output<string>;
    public readonly target!: pulumi.Output<string>;

    /**
     * Create a Link resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: LinkArgs, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.exists === undefined) && !opts.urn) {
                throw new Error("Missing required property 'exists'");
            }
            if ((!args || args.is_dir === undefined) && !opts.urn) {
                throw new Error("Missing required property 'is_dir'");
            }
            if ((!args || args.linked === undefined) && !opts.urn) {
                throw new Error("Missing required property 'linked'");
            }
            if ((!args || args.overwrite === undefined) && !opts.urn) {
                throw new Error("Missing required property 'overwrite'");
            }
            if ((!args || args.source === undefined) && !opts.urn) {
                throw new Error("Missing required property 'source'");
            }
            if ((!args || args.target === undefined) && !opts.urn) {
                throw new Error("Missing required property 'target'");
            }
            resourceInputs["exists"] = args ? args.exists : undefined;
            resourceInputs["is_dir"] = args ? args.is_dir : undefined;
            resourceInputs["linked"] = args ? args.linked : undefined;
            resourceInputs["overwrite"] = args ? args.overwrite : undefined;
            resourceInputs["source"] = args ? args.source : undefined;
            resourceInputs["target"] = args ? args.target : undefined;
            resourceInputs["result"] = undefined /*out*/;
        } else {
            resourceInputs["exists"] = undefined /*out*/;
            resourceInputs["is_dir"] = undefined /*out*/;
            resourceInputs["linked"] = undefined /*out*/;
            resourceInputs["overwrite"] = undefined /*out*/;
            resourceInputs["result"] = undefined /*out*/;
            resourceInputs["source"] = undefined /*out*/;
            resourceInputs["target"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(Link.__pulumiType, name, resourceInputs, opts);
    }
}

/**
 * The set of arguments for constructing a Link resource.
 */
export interface LinkArgs {
    exists: pulumi.Input<boolean>;
    is_dir: pulumi.Input<boolean>;
    linked: pulumi.Input<boolean>;
    overwrite: pulumi.Input<boolean>;
    source: pulumi.Input<string>;
    target: pulumi.Input<string>;
}
