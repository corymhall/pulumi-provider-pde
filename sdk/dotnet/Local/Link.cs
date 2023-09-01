// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Pde.Local
{
    [PdeResourceType("pde:local:Link")]
    public partial class Link : global::Pulumi.CustomResource
    {
        [Output("is_dir")]
        public Output<bool> Is_dir { get; private set; } = null!;

        [Output("linked")]
        public Output<bool> Linked { get; private set; } = null!;

        [Output("overwrite")]
        public Output<bool?> Overwrite { get; private set; } = null!;

        [Output("recursive")]
        public Output<bool?> Recursive { get; private set; } = null!;

        [Output("retain")]
        public Output<bool?> Retain { get; private set; } = null!;

        [Output("source")]
        public Output<string> Source { get; private set; } = null!;

        [Output("target")]
        public Output<string> Target { get; private set; } = null!;

        [Output("targets")]
        public Output<ImmutableArray<string>> Targets { get; private set; } = null!;


        /// <summary>
        /// Create a Link resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Link(string name, LinkArgs args, CustomResourceOptions? options = null)
            : base("pde:local:Link", name, args ?? new LinkArgs(), MakeResourceOptions(options, ""))
        {
        }

        private Link(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("pde:local:Link", name, null, MakeResourceOptions(options, id))
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
        /// Get an existing Link resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static Link Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new Link(name, id, options);
        }
    }

    public sealed class LinkArgs : global::Pulumi.ResourceArgs
    {
        [Input("overwrite")]
        public Input<bool>? Overwrite { get; set; }

        [Input("recursive")]
        public Input<bool>? Recursive { get; set; }

        [Input("retain")]
        public Input<bool>? Retain { get; set; }

        [Input("source", required: true)]
        public Input<string> Source { get; set; } = null!;

        [Input("target", required: true)]
        public Input<string> Target { get; set; } = null!;

        public LinkArgs()
        {
        }
        public static new LinkArgs Empty => new LinkArgs();
    }
}
