# Table: pokemon_item

An item is an object in the games which the player can pick up, keep in their bag,
and use in some manner. They have various uses, including healing, powering up,
helping catch Pokémon, or to access a new area.

## Examples

### Basic info

```sql
select
  name,
  id,
  cost,
  category,
  fling_power
from
  pokemon_item
```

### List all Item having cost more than 200 units

```sql
select
  name,
  id,
  cost
from
  pokemon_item
where
  cost >= 200
```
