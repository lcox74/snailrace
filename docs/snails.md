# Snails

Snails are the crux of the snailrace game, they are the entry point to the games
world. As an Role Playing Game (RPG), they need to feel like they belong and are
effected by things in the game. The snails in Snailrace must:

- Have their own breeds and how those breeds affect the stats and abilities
- Equip up to 5 items that modify the snails stats
- Breeding snails which combined the stats of the parent snail but still keep
  one of the parents breeds.
- Leveling and progressing, allowing the snail to improve its stats 
  progressively.

Each snail is built up of their own stats, which can be modified through breed,
equipment and leveling up. There are 4 main statistics that are used for racing,
each are in the range of `[1, 20]`:

- **Speed**: 
    This stat represents how fast the snail can go, its a relative value as the 
    snail's are limited to `5 units per step`.

- **Stamina**: 
    Stamina determines how long a snail can maintain its top speed before 
    getting fatigued.

- **Agility**: 
    Agility affects a snail's ability to navigate turns or obstacles on the 
    racecourse. Though there is no plans of adding obstacles to the racecourse,
    it can be seen as navigation through inperfections in the track.

- **Endurance**: 
    Endurance represents a snails's overall ability to withstand the rigors of 
    a race indicating how long a snail can race before becoming completely 
    exhausted.

The use of these stats can be used in the following function formula:

```go
stats := ... // Get Snail's Racing stats

// Initial Effective Speed
effective_speed := (5 * (stats.Speed / 20.0)) + rand.Float64()

// Calculate edurance factor and adjust the effective speed
endurance_factor := ((stats.Stamina + stats.Endurance) / 2.0) / 20.0
effective_speed *= endurance_factor

// Simulate agility by introducing a chance of losing speed
agility_chance := h.Agility / 20.0
if rand.Float64() < agility_chance {
    effective_speed *= 0.9 // Reduce speed by 10% if agility check fails
}

// Update Snail Position
position += effective_speed

```

Unlike the previous design of Starter, Amatuer, Expert snails. This spec acts 
more RPG like, where snails have a base stats and have a limited way to 
increase its stats. Snails can max out to **Level 5** at each level the snail
owner can increase a single stat. When the snails reach **Level 3, 4, 5** they
also get the option of selecting an **Ability** which are different for each
breed. These abilities modify the snail's stats futher and must be chosen 
selectively as the snail can only have **3 Abilities**.

Items are another way of increasing the stats of snails, these can be bought, 
found, traded or won. Items are ranked in tiers: `Common`, `Uncommon`, `Rare` 
and `Legendary`, each rank affect snails in a increasing order respectively and
also increase in difficulty to get. A snail can only hold **3 Items** at a time
and be equipped/dequipped between other snails. Each breed have a item 
requirement for breeding, which after breeding the item is destroyed.

## Snail Breeds

The following are the break down of the snails in the game world. This includes
the base states for each breed of snails, a description of their appearance,
some lore, unlockable abilities, and the lore on the requirements to breed.

### Lustrous Prismshell 🌟

| Stat      | Value |
| --------- | ----- |
| Agility   |   3   |
| Speed     |   4   |
| Endurance |   3   |
| Stamina   |   3   |

The Lustrous Prismshell is a captivating snail breed with a shell that resembles
a magnificent crystal prism. Its shell exhibits a stunning array of colors, 
reflecting and refracting light in a dazzling spectacle. The snail's body is 
sleek and graceful, perfectly complementing its mesmerizing exterior.

Lustrous Prismshells are found in ancient crystalline caves and hidden 
underground chambers. These snails thrive in the tranquil darkness, where the 
ambient glow of precious gemstones illuminates their surroundings. Their agility
allows them to navigate the intricate passages with ease, while their endurance
enables them to withstand the challenges of subterranean life.

The history of the Lustrous Prismshell snails is shrouded in ancient tales and 
mysticism. According to legend, these extraordinary snails were first discovered
by the ancient star gazers of Etheria. These stargazers believed that the 
Lustrous Prismshells were celestial messengers, embodying the secrets of the 
cosmos within their shimmering shells. It is said that the snails' iridescent 
hues were bestowed upon them by the moon goddess as a reward for their 
unwavering devotion to the night sky. Over the centuries, these snails have 
become revered symbols of celestial wisdom and enlightenment, their presence in
races serving as a reminder of the connection between the mortal realm and the 
infinite expanse of the universe.

**Skills**

- **Celestial Radiance**: Unleashes radiant energy that boosts the snail's 
  Speed (`+2`) and Agility (`+1`).
- **Stellar Leap**: Harnessing the power of the stars, the snail is able to
  perform swift and agile leaps, increasing its Agility (`+2`).
- **Luminary Aura**: Emitting a luminous aura, the snail enhances its 
  Speed (`+3`).
- **Radiant Shield**: Creating a shield of celestial light, the snail increases 
  its Endurance (`+2`) and gains protection against its preditors.
- **Solar Beam**: Focuses solar energy into a concentrated beam, unleashing a 
  powerful blast causing Stamina (`+4`) and Agility (`-1`). 
- **Astral Projection**: Projects an ethereal duplicate of the snail, confusing 
  opponents and causing them to make mistakes, resulting in reduced opponent 
  Agility (`+2`) and Endurance (`+2`).

**Breeding**

In the depths of the enchanted forest, where the starlight filters through the 
ancient trees, lies a hidden gem known as the **Crystal Shard**. This shimmering 
fragment, fallen from the celestial realm, possesses a captivating luminosity 
that resonates with the Lustrous Prismshell breed. When used in breeding, the 
Crystal Shard unlocks the snail's innate celestial powers, infusing its 
abilities with an otherworldly brilliance. As the Crystal Shard merges with the 
snail's essence, it bestows a unique ability modifier upon the breed, enhancing 
the celestial radiance it can emit. With the Crystal Shard's influence, the 
snail's Celestial Radiance ability gains an additional boost, increasing the 
`Speed +1`. This mystical union between the Lustrous Prismshell and the Crystal 
Shard unveils a snail of unparalleled grace and celestial prowess.

### Stormstrike Thunderhorn ⚡

| Stat      | Value |
| --------- | ----- |
| Agility   |   4   |
| Speed     |   4   |
| Endurance |   3   |
| Stamina   |   2   |

The Stormstrike Thunderhorn snail possesses a majestic shell that shimmers with 
electrifying energy. Its dark gray shell is adorned with intricate patterns 
reminiscent of storm clouds and lightning bolts. The snail's body is sleek and 
nimble, perfectly suited for swift movements across various terrains.

Stormstrike Thunderhorns thrive in regions of intense weather activity, such as 
stormy skies, lightning-laden landscapes, and lush rainforests. They are often 
found racing along winding paths amidst pouring rain, their agility enabling 
them to dodge falling drops with precision. These snails are at home in the 
heart of tempestuous weather.

The history of the Stormstrike Thunderhorn snails is intertwined with the 
tumultuous nature of the land of Electria. In ages past, the Thunderhorn snails 
were considered the chosen companions of the Storm Guardians, revered guardians 
of the realm who wielded the power of thunder and lightning. The snails were 
believed to be blessed by the gods of the tempest, granting them the ability to
harness the raw energy of storms. These powerful creatures were renowned for 
their participation in ancient races, where their lightning-quick speed and 
thunderous presence left spectators in awe. To this day, the Stormstrike 
Thunderhorns remain a symbol of power and resilience, embodying the fierce 
spirit of the ever-changing skies.

**Skills**

- **Thunderstorm Surge**: By summoning the power of thunderstorms, the snail 
  experiences an increase in Speed (`+2`) and Endurance (`+1`).
- **Thunderous Charge**: The snail charges forward with electrifying speed, 
  stunning opponents in its path. It gains increased Speed (`+3`) and 
  Endurance (`+1`).
- **Electroshield**: Creating an electric shield, the snail boosts its 
  Endurance (`+2`) and provides protection.
- **Lightning Strike**: Unleashing a powerful lightning strike, the snail 
  experiences a surge in Speed (`+2`) and Agility (`+2`).
- **Electrocharge**: Charges the snail's shell with electrical energy, 
  electrifying opponents on contact and reducing their Agility and Stamina, as 
  well as increasing Agility (`+2`) and Stamina (`+2`). 
- **Tempest Whirlwind**: Creates a powerful whirlwind around the snail, 
  deflecting incoming projectiles and providing an Agility (`+4`) boost to the 
  snail while also causing (`-1`).

**Breeding**

Deep within the heart of the stormy mountains, where thunder echoes and 
lightning dances across the sky, resides the fabled **Lightning Essence**. This 
electrifying substance is distilled from the very essence of thunderstorms, 
capturing the raw energy that courses through the atmosphere. The Stormstrike 
Thunderhorn breed is inexorably drawn to the power of lightning and requires 
the Lightning Essence to reach its full potential. When used in breeding, the 
Lightning Essence infuses the snail with intensified electrical prowess, 
elevating its thunderous abilities. The snail becomes electrified, crackling 
with energy and agility. The Lightning Essence brings forth a powerful ability 
modifier, empowering the Thunderstorm Surge with an increased boost. Both the 
`Speed +1` and `Endurance +1` enhancements are amplified, as the snail harnesses 
the raw power of thunderstorms to surge forward and leave its competitors 
stunned.

### Gilded Royalcrest 👑

| Stat      | Value |
| --------- | ----- |
| Agility   |   4   |
| Speed     |   3   |
| Endurance |   3   |
| Stamina   |   3   |

The Gilded Royalcrest snail is a majestic breed that possesses a shell adorned 
with intricate golden patterns, reminiscent of regal crowns. The shimmering gold
hue of their shells exudes an aura of nobility. Their bodies are robust and 
elegant, displaying the grace befitting their royal lineage.

Gilded Royalcrests are often found in enchanted forests with sunlit glades and 
meandering streams. These snails thrive in the presence of ancient trees and 
magical flora. They navigate the forest floor with endurance and elegance, their
golden shells blending harmoniously with the dappled sunlight that filters
through the canopy.

The history of the Gilded Royalcrest snails is intertwined with the ancient lore
of the Enchanted Kingdom. According to legend, these majestic creatures were 
once companions to the forest spirits that protected the mystical realms. The 
snails' golden shells were believed to be gifts from the ancient guardians, 
forged from the essence of the sun and infused with ancient magic. The 
Royalcrests served as guides to lost travelers, leading them through enchanted 
forests and hidden groves with their shimmering trails. As the Enchanted Kingdom
evolved over time, the snails became symbols of wisdom and nobility, their regal
appearance captivating all who beheld their resplendent beauty. Today, the 
Gilded Royalcrests stand as guardians of the natural world, their presence in 
races a testament to the enduring harmony between mortals and nature.

**Skills**

- **Regal Presence**: Exuding regal grace, the snail enhances its 
  Endurance (`+2`) and Stamina (`+1`).
- **Majestic Glide**: Gliding gracefully, the snail conserves Stamina and 
  enhances its Agility (`+2`).
- **Noble Resilience**: Tapping into inner resilience, the snail boosts its 
  Endurance (`+3`).
- **Royal Command**: Issuing a regal command, the snail temporarily boosts the
  Speed (`+2`) and Stamina (`+2`) of itself and its allies.
- **Regal Enchantment**: Enchants the snail's shell with regal energy, boosting 
  the snail's Endurance (`+2`) and Stamina (`+2`) while decreasing the 
  opponent's Speed. 
- **Golden Sheen**: The snail's shell is coated a golden sheen allowing it to
  slip through the wind increasing Speed (`+3`).

**Breeding**

In the realm of elegance and regality, where opulence and grace intertwine, lies
the coveted **Regal Crest**. Crafted with intricate artistry and adorned with 
precious gems, this emblem symbolizes the nobility and grandeur of the Gilded 
Royalcrest breed. To breed a snail of such majestic lineage, one must possess 
the revered Regal Crest. When united with the snail, the Regal Crest unlocks a 
regal presence within the breed, enhancing its endurance and stamina. As the 
Regal Crest melds with the snail's essence, it imparts an ability modifier that
further elevates its inherent regal qualities. The snail's ability to command 
respect and admiration reaches new heights as the Royal Crest bestows an 
additional boost to both `Endurance +1` and `Stamina +1` when activating the 
Regal Presence ability. With the influence of the Regal Crest, the snail becomes
an epitome of unwavering resilience, commanding attention on the racetrack with
its noble bearing.

### Bioengineered Circuitshell 🤖

| Stat      | Value |
| --------- | ----- |
| Agility   |   3   |
| Speed     |   2   |
| Endurance |   4   |
| Stamina   |   4   |

The Bioengineered Circuitshell is a unique snail breed, created through advanced
scientific methods. Its shell is a combination of sleek metal plating and 
translucent bioengineered materials. Embedded within the shell are intricate 
circuits that pulse with energy. The snail's body exhibits a streamlined design
optimized for speed and efficiency.

Bioengineered Circuitshells are adaptable to various environments, but they are 
often found in technologically advanced laboratories and controlled habitats. 
These snails are the result of scientific experimentation, designed to excel in
races that demand both Endurance and stamina. Their metallic shells and 
bioengineered enhancements make them stand out among other snail breeds.

The history of the Bioengineered Circuitshell snails is a tale of innovation and
scientific achievement. Created by the visionary technomages of the ancient 
Technomage Guild, these snails represent a revolutionary leap in snail racing.
Faced with the limitations of natural evolution, the technomages embarked on a
daring experiment to create a breed that would surpass the boundaries of natural
capabilities. Through their relentless pursuit of knowledge, they combined 
organic matter with advanced technology, resulting in the birth of the 
Circuitshells. These artificial wonders quickly made a name for themselves in 
the racing circuit, their unmatched speed and stamina capturing the attention of
enthusiasts and researchers alike. The Bioengineered Circuitshells have become 
living testaments to the boundless possibilities of science and its impact on 
the world of snail racing.

**Skills**

- **Technomagical Enhancement**: Upgrading through technomagical enhancements, 
  the snail improves its Speed (`+2`) and Stamina (`+1`).
- **Cyber Dash**: Activating thrusters for a burst of incredible Speed, the 
  snail accelerates with cybernetic precision, gaining increased Speed (`+3`).
- **Nano-Repair Matrix**: Deploying a nano-repair matrix, the snail boosts its 
  Stamina (`+2`) and rapidly recovers health.
- **Overclocked Surge**: Pushing the limits with an overclocked surge, the snail
  temporarily increases its Speed (`+2`) and Agility (`+2`) to extraordinary
  levels.
- **Cybernetic Overdrive**: Activates a temporary overdrive mode, greatly 
  boosting the snail's Speed (`+4`) but causing a slight decrease in 
  Endurance (`-1`).
- **technomagnetic Disruption**: Emits disruptive technomagnetic waves, 
  temporarily disabling opponent abilities giving the snail improved 
  Endurance (`+2`) and Agility (`+2`)


**Breeding**

In the realm of technological marvels, where gears turn and circuits hum with 
energy, lies the enigmatic **Techno Core**. Crafted by master engineers and 
infused with cutting-edge advancements, this cybernetic implant holds the key to
unlocking the full potential of the Bioengineered Circuitshell breed. To breed a
snail of such technological prowess, one must possess the elusive Techno Core. 
When integrated into the snail's neural framework, the Techno Core augments its
inherent abilities, elevating its speed and stamina to unparalleled levels. As
the Techno Core merges with the snail's essence, it imparts a unique ability 
modifier that further enhances its technological enhancements. The snail's 
connection to the digital realm deepens, unlocking an amplified boost within the
Technomagical Enhancement ability. Both the `Speed +1` and `Stamina +1` 
enhancements are heightened, propelling the snail forward with cybernetic 
precision and unprecedented endurance.

### Emberwing Infernoshell 🔥

| Stat      | Value |
| --------- | ----- |
| Agility   |   4   |
| Speed     |   2   |
| Endurance |   3   |
| Stamina   |   3   |

The Emberwing Infernoshell snail boasts a fiery red shell adorned with intricate
patterns reminiscent of swirling flames. Its vibrant hues radiate heat, creating
an aura of intensity. The snail's body is sleek and agile, reflecting its 
ability to maneuver through fiery environments.

Emberwing Infernoshells prefer habitats of intense heat, such as volcanic 
regions, scorching deserts, and lava-filled landscapes. They navigate these 
treacherous terrains with impressive endurance and agility, their fiery shells 
providing resistance to extreme temperatures. These snails leave behind trails 
of scorched earth as they blaze through the most challenging racecourses.

The history of the Emberwing Infernoshell snails is steeped in the fiery depths 
of the Fiery Abyss. According to ancient folklore, the Infernoshells were born 
from the essence of the elemental flames that surged through the volcanic realm.
These snails were said to be chosen by the Fire Elemental Lords, who infused 
their shells with the power of the inferno. Tales speak of great races held 
amidst the lava flows, where the Emberwing Infernoshells raced alongside 
erupting volcanoes, leaving behind trails of scorching heat. The snails were
hailed as symbols of passion and determination, their presence evoking the 
primal forces of fire. To this day, the Emberwing Infernoshells are revered as 
emissaries of the flames, their history intertwined with the ever-burning 
essence of the Fiery Abyss.

**Skills**

- **Inferno Fury**: Channeling the fires of the inferno, the snail increases its
  Endurance (`+2`) and Agility (`+1`).
- **Blazing Dash**: Igniting with blazing speed, the snail leaves trails of fire
  and increases its Agility (`+3`).
- **Magma Shield**: Summoning a shield of molten magma, the snail boosts its 
  Endurance (`+3`) and gains heat resistance.
- **Infernal Rage**: Unleashing an infernal rage, the snail boosts its 
  Endurance (`+2`) and Stamina (`+2`) to extraordinary levels.
- **Inferno Blaze**: Ignites the snail's shell with an intense inferno, greatly 
  increasing Speed (`+4`) and Agility (`+2`) but greatly decreasing 
  Samina (`-2`).
- **Molten Shell**: Hardens the snail's shell to molten levels, providing 
  increased Stamina (`+2`) and Endurance (`+1`) and resistance to opponent 
  attacks.

**Breeding**

Within the molten depths of volcanic realms, where rivers of lava cascade and 
flames dance in a mesmerizing inferno, lies the coveted **Inferno Ember**. This
smoldering fragment, born from the heart of a fiery volcano, possesses an 
intense and relentless energy that resonates with the Emberwing Infernoshell 
breed. To breed a snail of such blazing determination, one must possess the 
scorching Inferno Ember. As the Ember fuses with the snail's essence, it engulfs
the breed in an infernal blaze, igniting its fiery powers and instilling an 
unrivaled passion within. The Inferno Ember imparts an ability modifier that 
amplifies the snail's innate abilities. The Inferno Fury ability is intensified,
augmenting both the `Endurance +1` and `Agility +1` boosts. With the influence 
of the Inferno Ember, the snail becomes an unstoppable force, channeling the 
fires of the inferno and leaving trails of blazing glory on the racetrack.