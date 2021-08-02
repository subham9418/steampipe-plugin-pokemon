package pokemon

import (
	"context"

	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tablePokemonItem(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pokemon_item",
		Description: "Pokémon are the creatures that inhabit the world of the Pokémon games.",
		List: &plugin.ListConfig{
			Hydrate: listItem,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name"}),
			// TODO: Add support for 'id' key column
			//KeyColumns: plugin.AnyColumn([]string{"id", "name"}),
			Hydrate: getItem,
			// Bad error message is a result of https://github.com/mtslzr/pokeapi-go/issues/29
			ShouldIgnoreError: isNotFoundError([]string{"invalid character 'N' looking for beginning of value"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cost",
				Description: "The price of this item in stores.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getItem,
			},
			{
				Name:        "fling_power",
				Description: "The power of the move Fling when used with this item.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getItem,
			},
			{
				Name:        "fling_effect",
				Description: "The effect of the move Fling when used with this item.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "attributes",
				Description: "A list of attributes this item has.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "category",
				Description: "The category of items this item falls into.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "effect_entries",
				Description: "The effect of this ability listed in different languages.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "flavor_text_entries",
				Description: "The flavor text of this ability listed in different languages.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "game_indices",
				Description: "A list of game indices relevent to this item by generation.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "sprites",
				Description: "A set of sprites used to depict this item in the game.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "held_by_pokemon",
				Description: "A list of Pokémon that might be found in the wild holding this item.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "baby_trigger_for",
				Description: "An evolution chain this item requires to produce a bay during mating.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "machines",
				Description: "A list of the machines related to this item.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getItem,
			},
			{
				Name:        "id",
				Description: "The identifier for this resource.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getItem,
				Transform:   transform.FromGo(),
			},
			// Standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

func listItem(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listItem")

	offset := 0

	for true {
		resources, err := pokeapi.Resource("item", offset)

		if err != nil {
			plugin.Logger(ctx).Error("pokemon_item.listItem", "query_error", err)
			return nil, err
		}

		for _, item := range resources.Results {
			d.StreamListItem(ctx, item)
		}

		// No next URL returned
		if len(resources.Next) == 0 {
			break
		}

		urlOffset, err := extractUrlOffset(resources.Next)
		if err != nil {
			plugin.Logger(ctx).Error("pokemon_item.listItem", "extract_url_offset_error", err)
			return nil, err
		}

		// Set next offset
		offset = urlOffset
	}

	return nil, nil
}

func getItem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getItem")

	var name string

	if h.Item != nil {
		result := h.Item.(structs.Result)
		name = result.Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	logger.Debug("Name", name)

	pokemon, err := pokeapi.Item(name)

	if err != nil {
		plugin.Logger(ctx).Error("pokemon_item.itemGet", "query_error", err)
		return nil, err
	}

	return pokemon, nil
}
