import type { Recipe } from "@/lib/types";

export const recipes: Recipe[] = [
  {
    id: "1",
    slug: "classic-carbonara",
    title: { en: "Classic Carbonara", ru: "Классическая карбонара" },
    description: {
      en: "Silky Roman pasta with guanciale, egg yolks, and Pecorino Romano. A timeless comfort dish that comes together in minutes.",
      ru: "Нежная римская паста с гуанчале, яичными желтками и пекорино романо. Классическое блюдо, готовится за считанные минуты."
    },
    image:
      "https://images.unsplash.com/photo-1612874742237-6526221588e3?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Plate of spaghetti carbonara with crispy guanciale",
      ru: "Тарелка спагетти карбонара с хрустящим гуанчале"
    },
    cuisine: "italian",
    category: "dinner",
    mealType: ["dinner"],
    diet: [],
    difficulty: "medium",
    prepTimeMinutes: 10,
    cookTimeMinutes: 15,
    totalTimeMinutes: 25,
    servings: 4,
    rating: 4.8,
    ingredients: {
      en: [
        "400 g spaghetti",
        "200 g guanciale, diced",
        "4 large egg yolks",
        "100 g Pecorino Romano, finely grated",
        "Freshly cracked black pepper",
        "Salt for pasta water"
      ],
      ru: [
        "400 г спагетти",
        "200 г гуанчале, нарезать кубиками",
        "4 крупных яичных желтка",
        "100 г пекорино романо, мелко натёртого",
        "Свежемолотый чёрный перец",
        "Соль для воды пасты"
      ]
    },
    steps: {
      en: [
        "Bring a large pot of salted water to a boil and cook the spaghetti until al dente.",
        "Meanwhile, fry the diced guanciale in a dry pan over medium heat until golden and crisp.",
        "Whisk the egg yolks with the grated Pecorino and a generous pinch of black pepper.",
        "Reserve a cup of pasta water, then drain the spaghetti and add it to the pan with the guanciale.",
        "Remove from heat and immediately toss with the egg mixture, adding pasta water until creamy.",
        "Serve at once with extra Pecorino and black pepper."
      ],
      ru: [
        "Вскипятите в большой кастрюле подсоленную воду и отварите спагетти аль денте.",
        "Тем временем обжарьте нарезанный гуанчале на сухой сковороде на среднем огне до золотистости.",
        "Взбейте яичные желтки с натёртым пекорино и щедрой щепоткой чёрного перца.",
        "Отлейте стакан воды от пасты, слейте спагетти и добавьте их на сковороду с гуанчале.",
        "Снимите с огня и сразу перемешайте с яичной смесью, добавляя воду от пасты до кремовой консистенции.",
        "Подавайте сразу же с дополнительным пекорино и чёрным перцем."
      ]
    },
    nutrition: { calories: 580, protein: 24, fat: 28, carbs: 60 },
    tips: {
      en: [
        "Keep the pan off the heat when adding eggs to prevent scrambling.",
        "Use Pecorino Romano, not Parmesan, for an authentic Roman flavor."
      ],
      ru: [
        "Снимайте сковороду с огня перед добавлением яиц, чтобы они не свернулись.",
        "Используйте пекорино романо, а не пармезан, для аутентичного римского вкуса."
      ]
    },
    tags: ["pasta", "comfort-food", "italian", "quick"],
    datePublished: "2024-01-15"
  },
  {
    id: "2",
    slug: "margherita-pizza",
    title: { en: "Margherita Pizza", ru: "Пицца Маргарита" },
    description: {
      en: "Neapolitan classic with San Marzano tomato, fresh mozzarella, and basil on a blistered wood-fired crust.",
      ru: "Неаполитанская классика с томатами сан-марцано, свежей моцареллой и базиликом на blistered-корже."
    },
    image:
      "https://images.unsplash.com/photo-1574071318508-1cdbab80d002?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Margherita pizza with basil leaves on a wooden board",
      ru: "Пицца Маргарита с листьями базилика на деревянной доске"
    },
    cuisine: "italian",
    category: "baking",
    mealType: ["dinner"],
    diet: ["vegetarian"],
    difficulty: "hard",
    prepTimeMinutes: 20,
    cookTimeMinutes: 12,
    totalTimeMinutes: 32,
    servings: 4,
    rating: 4.7,
    ingredients: {
      en: [
        "500 g tipo 00 flour",
        "325 ml warm water",
        "10 g salt",
        "5 g instant yeast",
        "400 g San Marzano tomatoes, crushed",
        "250 g fresh mozzarella, torn",
        "Fresh basil leaves",
        "Extra virgin olive oil"
      ],
      ru: [
        "500 г муки тип 00",
        "325 мл тёплой воды",
        "10 г соли",
        "5 г быстродействующих дрожжей",
        "400 г томатов сан-марцано, размятых",
        "250 г свежей моцареллы, порванной руками",
        "Свежие листья базилика",
        "Оливковое масло extra virgin"
      ]
    },
    steps: {
      en: [
        "Mix the flour, water, salt, and yeast into a shaggy dough and knead for 10 minutes until smooth.",
        "Cover and let the dough rise at room temperature for 8 hours, or refrigerate overnight.",
        "Divide into four balls and rest for another 2 hours.",
        "Stretch each ball into a thin round on a floured surface.",
        "Spread crushed tomatoes over the dough and add torn mozzarella.",
        "Bake in the hottest oven possible (250 °C+) for 7–10 minutes until the crust is blistered.",
        "Finish with fresh basil and a drizzle of olive oil."
      ],
      ru: [
        "Замесите из муки, воды, соли и дрожжей тесто и вымешивайте 10 минут до гладкости.",
        "Накройте и дайте тесту подойти при комнатной температуре 8 часов или уберите в холодильник на ночь.",
        "Разделите на четыре шарика и дайте постоять ещё 2 часа.",
        "Растяните каждый шарик в тонкий круг на посыпанной мукой поверхности.",
        "Намажьте размятые томаты на тесто и добавьте порванную моцареллу.",
        "Выпекайте в максимально горячей духовке (250 °C и выше) 7–10 минут до пузырей на корже.",
        "Украсьте свежим базиликом и сбрызните оливковым маслом."
      ]
    },
    nutrition: { calories: 410, protein: 18, fat: 16, carbs: 50 },
    tips: {
      en: [
        "Use a pizza stone or steel preheated for at least 45 minutes for a crisp base.",
        "Don't overload with toppings — less is more for Margherita."
      ],
      ru: [
        "Используйте камень или сталь для пиццы, разогретые не менее 45 минут, для хрустящей основы.",
        "Не перегружайте начинкой — для Маргариты меньше значит лучше."
      ]
    },
    tags: ["pizza", "baking", "italian", "vegetarian"],
    datePublished: "2024-02-20"
  },
  {
    id: "3",
    slug: "georgian-khachapuri",
    title: { en: "Georgian Khachapuri", ru: "Грузинский хачапури" },
    description: {
      en: "Boat-shaped cheese-filled bread topped with a runny egg and a knob of butter — Georgia's most iconic comfort food.",
      ru: "Лодочка из дрожжевого хлеба с сырной начинкой, увенчанная жидким яйцом и кусочком масла — самое культовое грузинское блюдо."
    },
    image:
      "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Khachapuri boat bread with egg yolk and butter",
      ru: "Хачапури-лодочка с яичным желтком и маслом"
    },
    cuisine: "georgian",
    category: "baking",
    mealType: ["lunch", "dinner"],
    diet: ["vegetarian"],
    difficulty: "medium",
    prepTimeMinutes: 20,
    cookTimeMinutes: 20,
    totalTimeMinutes: 40,
    servings: 4,
    rating: 4.9,
    ingredients: {
      en: [
        "500 g all-purpose flour",
        "300 ml warm milk",
        "7 g instant yeast",
        "1 tsp sugar",
        "1 tsp salt",
        "400 g sulguni cheese, grated",
        "100 g mozzarella, grated",
        "4 eggs",
        "50 g butter"
      ],
      ru: [
        "500 г пшеничной муки",
        "300 мл тёплого молока",
        "7 г быстродействующих дрожжей",
        "1 ч. л. сахара",
        "1 ч. л. соли",
        "400 г сыра сулугуни, натёртого",
        "100 г моцареллы, натёртой",
        "4 яйца",
        "50 г сливочного масла"
      ]
    },
    steps: {
      en: [
        "Mix flour, warm milk, yeast, sugar, and salt into a soft dough and knead for 8 minutes.",
        "Cover and let rise in a warm spot for about 1.5 hours until doubled.",
        "Combine the grated sulguni and mozzarella in a bowl.",
        "Divide the dough into four portions and roll each into an oval.",
        "Pinch the ends to form a boat shape and fill the centre with the cheese mixture.",
        "Bake at 220 °C for 12 minutes, then crack an egg into each boat and bake 3 more minutes.",
        "Add a knob of butter on top and serve immediately, tearing the crust to stir into the egg."
      ],
      ru: [
        "Замесите из муки, тёплого молока, дрожжей, сахара и соли мягкое тесто и вымешивайте 8 минут.",
        "Накройте и дайте подойти в тёплом месте около 1,5 часов до удвоения объёма.",
        "Соедините натёртый сулугуни и моцареллу в миске.",
        "Разделите тесто на четыре части и раскатайте каждую в овал.",
        "Защипните края, формируя лодочку, и наполните центр сырной смесью.",
        "Выпекайте при 220 °C 12 минут, затем вбейте в каждую лодочку яйцо и пеките ещё 3 минуты.",
        "Положите сверху кусочек масла и подавайте сразу же, отрывая корж и макая в яйцо."
      ]
    },
    nutrition: { calories: 540, protein: 24, fat: 26, carbs: 48 },
    tips: {
      en: [
        "If you can't find sulguni, mix feta and mozzarella for a similar stretch.",
        "Serve right out of the oven for the best runny-egg effect."
      ],
      ru: [
        "Если не нашли сулугуни, смешайте фету и моцареллу для похожей тягучести.",
        "Подавайте прямо из духовки ради идеального жидкого яйца."
      ]
    },
    tags: ["bread", "cheese", "georgian", "comfort-food"],
    datePublished: "2024-03-10"
  },
  {
    id: "4",
    slug: "chicken-ramen",
    title: { en: "Chicken Ramen", ru: "Куриный рамен" },
    description: {
      en: "A soul-warming bowl with rich chicken broth, tender marinated chicken, soft-boiled egg, and springy noodles.",
      ru: "Согревающая тарелка с насыщенным куриным бульоном, нежной маринованной курицей, яйцом пашот и упругими лапшой."
    },
    image:
      "https://images.unsplash.com/photo-1569718212165-3a8278d5f624?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Bowl of chicken ramen with egg and green onions",
      ru: "Тарелка куриного рамена с яйцом и зелёным луком"
    },
    cuisine: "japanese",
    category: "dinner",
    mealType: ["dinner"],
    diet: [],
    difficulty: "hard",
    prepTimeMinutes: 25,
    cookTimeMinutes: 45,
    totalTimeMinutes: 70,
    servings: 4,
    rating: 4.6,
    ingredients: {
      en: [
        "1.5 L chicken stock",
        "2 chicken thighs, boneless",
        "4 portions fresh ramen noodles",
        "4 eggs",
        "3 tbsp soy sauce",
        "2 tbsp mirin",
        "1 tbsp miso paste",
        "2 spring onions, sliced",
        "1 sheet nori, cut into strips"
      ],
      ru: [
        "1,5 л куриного бульона",
        "2 куриных бедра, без кости",
        "4 порции свежей рамен-лапши",
        "4 яйца",
        "3 ст. л. соевого соуса",
        "2 ст. л. мирина",
        "1 ст. л. пасты мисо",
        "2 стебля зелёного лука, нарезанного",
        "1 лист нори, нарезать полосками"
      ]
    },
    steps: {
      en: [
        "Simmer the chicken stock for 20 minutes and keep warm.",
        "Marinate the chicken thighs in 2 tbsp soy sauce and mirin for 15 minutes, then pan-sear until cooked through.",
        "Soft-boil the eggs for 6.5 minutes, then peel and marinate in remaining soy sauce.",
        "Whisk the miso paste into the warm stock until dissolved.",
        "Cook the ramen noodles according to package directions and drain.",
        "Divide noodles into bowls, ladle over the hot misa stock, and slice the chicken on top.",
        "Finish with halved eggs, spring onion, and nori strips."
      ],
      ru: [
        "Варите куриный бульон 20 минут и держите тёплым.",
        "Замаринуйте куриные бёдра в 2 ст. л. соевого соуса и мирина на 15 минут, затем обжарьте до готовности.",
        "Сварите яйца всмятку 6,5 минуты, очистите и замаринуйте в оставшемся соевом соусе.",
        "Взбейте пасту мисо в тёплом бульоне до растворения.",
        "Отварите рамен-лапшу согласно инструкции на упаковке и слейте воду.",
        "Разложите лапшу по тарелкам, налейте горячий бульон с мисо и выложите нарезанную курицу сверху.",
        "Украсьте половинками яиц, зелёным луком и полосками нори."
      ]
    },
    nutrition: { calories: 480, protein: 30, fat: 18, carbs: 45 },
    tips: {
      en: [
        "For a richer broth, simmer chicken wings with ginger and garlic for 2 hours.",
        "Don't overcook the eggs — 6.5 minutes gives the perfect jammy yolk."
      ],
      ru: [
        "Для более насыщенного бульона варите куриные крылья с имбирём и чесноком 2 часа.",
        "Не переваривайте яйца — 6,5 минуты дают идеальный жидкий желток."
      ]
    },
    tags: ["noodles", "japanese", "soup", "comfort-food"],
    datePublished: "2024-04-05"
  },
  {
    id: "5",
    slug: "vegetable-curry",
    title: { en: "Vegetable Curry", ru: "Овощное карри" },
    description: {
      en: "Fragrant Indian-style curry with tender vegetables in a coconut-tomato sauce, ready in under an hour.",
      ru: "Ароматное индийское карри с нежными овощами в кокосово-томатном соусе, готовится меньше часа."
    },
    image:
      "https://images.unsplash.com/photo-1455619452474-d2be8b1e70cd?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Bowl of vegetable curry with rice and cilantro",
      ru: "Тарелка овощного карри с рисом и кинзой"
    },
    cuisine: "indian",
    category: "dinner",
    mealType: ["dinner"],
    diet: ["vegan", "gluten-free"],
    difficulty: "medium",
    prepTimeMinutes: 15,
    cookTimeMinutes: 30,
    totalTimeMinutes: 45,
    servings: 4,
    rating: 4.5,
    ingredients: {
      en: [
        "1 tbsp coconut oil",
        "1 onion, diced",
        "3 garlic cloves, minced",
        "1 tbsp grated ginger",
        "2 tbsp curry powder",
        "400 ml coconut milk",
        "200 g tomatoes, chopped",
        "300 g mixed vegetables (cauliflower, carrots, peas)",
        "Fresh cilantro to serve"
      ],
      ru: [
        "1 ст. л. кокосового масла",
        "1 луковица, нарезать кубиками",
        "3 зубчика чеснока, измельчить",
        "1 ст. л. натёртого имбиря",
        "2 ст. л. порошка карри",
        "400 мл кокосового молока",
        "200 г томатов, нарезанных",
        "300 г овощной смеси (цветная капуста, морковь, горошек)",
        "Свежая кинза для подачи"
      ]
    },
    steps: {
      en: [
        "Heat the coconut oil in a large pan and sauté the diced onion until soft.",
        "Add the garlic, ginger, and curry powder and stir for a minute until fragrant.",
        "Pour in the coconut milk and tomatoes and bring to a gentle simmer.",
        "Add the mixed vegetables, cover, and cook for 20 minutes until tender.",
        "Season to taste and serve hot with rice and fresh cilantro."
      ],
      ru: [
        "Разогрейте кокосовое масло в большой сковороде и обжарьте нарезанный лук до мягкости.",
        "Добавьте чеснок, имбирь и порошок карри, перемешивайте минуту до аромата.",
        "Влейте кокосовое молоко и томаты, доведите до лёгкого кипения.",
        "Добавьте овощи, накройте крышкой и тушите 20 минут до мягкости.",
        "Посолите по вкусу и подавайте горячим с рисом и свежей кинзой."
      ]
    },
    nutrition: { calories: 320, protein: 10, fat: 14, carbs: 38 },
    tips: {
      en: [
        "Toast the curry powder in the oil for 30 seconds to wake up the spices.",
        "Swap the vegetables for whatever is seasonal — pumpkin and green beans work well."
      ],
      ru: [
        "Прокалите порошок карри в масле 30 секунд, чтобы раскрыть специи.",
        "Замените овощи на сезонные — хорошо подходят тыква и стручковая фасоль."
      ]
    },
    tags: ["curry", "indian", "vegan", "gluten-free", "comfort-food"],
    datePublished: "2024-05-12"
  },
  {
    id: "6",
    slug: "caesar-salad",
    title: { en: "Caesar Salad", ru: "Салат Цезарь" },
    description: {
      en: "Crisp romaine leaves tossed in a creamy anchovy dressing with crunchy croutons and shaved Parmesan.",
      ru: "Хрустящие листья романо в кремовом анчоусном соусе с сухариками и тонтко нарезанным пармезаном."
    },
    image:
      "https://images.unsplash.com/photo-1546793665-c74683f339c1?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Caesar salad with croutons and parmesan on a plate",
      ru: "Салат Цезарь с сухариками и пармезаном на тарелке"
    },
    cuisine: "american",
    category: "salads",
    mealType: ["lunch"],
    diet: [],
    difficulty: "easy",
    prepTimeMinutes: 15,
    cookTimeMinutes: 5,
    totalTimeMinutes: 20,
    servings: 2,
    rating: 4.4,
    ingredients: {
      en: [
        "2 heads romaine lettuce, chopped",
        "2 slices stale bread, cubed",
        "3 anchovy fillets, chopped",
        "1 garlic clove",
        "2 egg yolks",
        "1 tbsp Dijon mustard",
        "50 ml olive oil",
        "30 g Parmesan, shaved",
        "Lemon juice to taste"
      ],
      ru: [
        "2 кочана салата романо, нарезать",
        "2 ломтика черствого хлеба, кубиками",
        "3 филе анчоусов, нарезать",
        "1 зубчик чеснока",
        "2 яичных желтка",
        "1 ст. л. дижонской горчицы",
        "50 мл оливкового масла",
        "30 г пармезана, тонтко нарезать",
        "Лимонный сок по вкусу"
      ]
    },
    steps: {
      en: [
        "Toast the bread cubes in a little olive oil until golden and set aside.",
        "Mash the garlic and anchovies together with a pinch of salt.",
        "Whisk in the egg yolks and Dijon, then slowly stream in the olive oil to form a thick dressing.",
        "Add a squeeze of lemon juice and toss the romaine with the dressing.",
        "Top with croutons and shaved Parmesan and serve at once."
      ],
      ru: [
        "Поджарьте хлебные кубики на небольшом количестве оливкового масла до золотистости и отложите.",
        "Разомните чеснок и анчоусы вместе со щепоткой соли.",
        "Вбейте яичные желтки и горчицу, затем медленно вливайте оливковое масло, взбивая до густого соуса.",
        "Добавьте лимонный сок и перемешайте романо с соусом.",
        "Посыпьте сухариками и пармезаном и подавайте сразу."
      ]
    },
    nutrition: { calories: 280, protein: 8, fat: 22, carbs: 12 },
    tips: {
      en: [
        "Use a ceramic bowl to mash the anchovies — it keeps the flavour bright.",
        "Dress the leaves just before serving so the croutons stay crisp."
      ],
      ru: [
        "Используйте керамическую миску для анчоусов — это сохраняет яркость вкуса.",
        "Заправляйте листья прямо перед подачей, чтобы сухарики оставались хрустящими."
      ]
    },
    tags: ["salad", "american", "lunch", "classic"],
    datePublished: "2024-05-30"
  },
  {
    id: "7",
    slug: "borscht",
    title: { en: "Borscht", ru: "Борщ" },
    description: {
      en: "Deep ruby Ukrainian beet soup with tender cabbage, potatoes, and a bright spoonful of smetana.",
      ru: "Насыщенный рубиновый украинский свекольный суп с нежной капустой, картофелем и яркой ложкой сметаны."
    },
    image:
      "https://images.unsplash.com/photo-1547592180-85f173990554?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Bowl of borscht with sour cream and dill",
      ru: "Тарелка борща со сметаной и укропом"
    },
    cuisine: "ukrainian",
    category: "soups",
    mealType: ["lunch", "dinner"],
    diet: ["vegetarian"],
    difficulty: "medium",
    prepTimeMinutes: 20,
    cookTimeMinutes: 40,
    totalTimeMinutes: 60,
    servings: 6,
    rating: 4.7,
    ingredients: {
      en: [
        "2 medium beets, grated",
        "1 onion, diced",
        "2 carrots, sliced",
        "1/4 white cabbage, shredded",
        "3 potatoes, cubed",
        "1 L vegetable stock",
        "2 tbsp tomato paste",
        "1 tbsp lemon juice",
        "Smetana or sour cream to serve",
        "Fresh dill"
      ],
      ru: [
        "2 средние свёклы, натереть",
        "1 луковица, нарезать кубиками",
        "2 моркови, нарезать",
        "1/4 кочана белокочанной капусты, нашинковать",
        "3 картофелины, кубиками",
        "1 л овощного бульона",
        "2 ст. л. томатной пасты",
        "1 ст. л. лимонного сока",
        "Сметана для подачи",
        "Свежий укроп"
      ]
    },
    steps: {
      en: [
        "Sauté the onion and carrots in a little oil until soft.",
        "Add the grated beets and tomato paste and cook for 5 minutes.",
        "Pour in the vegetable stock and bring to a boil.",
        "Add the potatoes and cabbage, then simmer for 25 minutes until tender.",
        "Stir in the lemon juice to keep the colour bright.",
        "Serve hot with a spoonful of smetana and fresh dill."
      ],
      ru: [
        "Обжарьте лук и морковь на небольшом количестве масла до мягкости.",
        "Добавьте натёртую свёклу и томатную пасту, готовьте 5 минут.",
        "Влейте овощной бульон и доведите до кипения.",
        "Добавьте картофель и капусту, варите 25 минут до мягкости.",
        "Влейте лимонный сок, чтобы сохранить яркий цвет.",
        "Подавайте горячим с ложкой сметаны и свежим укропом."
      ]
    },
    nutrition: { calories: 220, protein: 6, fat: 7, carbs: 28 },
    tips: {
      en: [
        "A splash of lemon juice keeps the beets ruby-red rather than fading.",
        "Borscht is even better the next day once the flavours have married."
      ],
      ru: [
        "Ложка лимонного сока сохраняет свёкле рубиновый цвет, не давая ей выцести.",
        "Борщ на следующий день вкуснее, когда вкусы «соединятся»."
      ]
    },
    tags: ["soup", "ukrainian", "beetroot", "vegetarian", "classic"],
    datePublished: "2024-06-18"
  },
  {
    id: "8",
    slug: "shakshuka",
    title: { en: "Shakshuka", ru: "Шакшука" },
    description: {
      en: "Eggs poached in a smoky tomato and pepper sauce with cumin and paprika — a classic Middle Eastern breakfast.",
      ru: "Яйца-пашот в дымчатом томатно-перечном соусе с кумином и паприкой — классический завтрак Ближнего Востока."
    },
    image:
      "https://images.unsplash.com/photo-1590412200988-a436970781fa?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Shakshuka eggs in tomato sauce in a cast iron pan",
      ru: "Шакшука — яйца в томатном соусе в чугунной сковороде"
    },
    cuisine: "mediterranean",
    category: "breakfast",
    mealType: ["breakfast"],
    diet: ["vegetarian", "gluten-free"],
    difficulty: "easy",
    prepTimeMinutes: 10,
    cookTimeMinutes: 20,
    totalTimeMinutes: 30,
    servings: 4,
    rating: 4.6,
    ingredients: {
      en: [
        "1 tbsp olive oil",
        "1 onion, diced",
        "1 red bell pepper, sliced",
        "3 garlic cloves, minced",
        "1 tsp cumin",
        "1 tsp paprika",
        "400 g chopped tomatoes",
        "4 eggs",
        "Salt and pepper to taste",
        "Fresh parsley"
      ],
      ru: [
        "1 ст. л. оливкового масла",
        "1 луковица, нарезать кубиками",
        "1 красный болгарский перец, нарезать",
        "3 зубчика чеснока, измельчить",
        "1 ч. л. кумина",
        "1 ч. л. паприки",
        "400 г нарезанных томатов",
        "4 яйца",
        "Соль и перец по вкусу",
        "Свежая петрушка"
      ]
    },
    steps: {
      en: [
        "Warm the olive oil in a skillet and sauté the onion and pepper until soft.",
        "Stir in the garlic, cumin, and paprika and cook for a minute.",
        "Add the tomatoes and simmer for 10 minutes until thickened.",
        "Make four wells in the sauce and crack an egg into each.",
        "Cover and cook for 4 to 6 minutes until the whites are set but yolks remain runny.",
        "Season and scatter with fresh parsley before serving."
      ],
      ru: [
        "Разогрейте оливковое масло в сковороде и обжарьте лук и перец до мягкости.",
        "Добавьте чеснок, кумин и паприку, готовьте минуту.",
        "Добавьте томаты, тушите 10 минут до загустения.",
        "Сделайте в соусе четыре улубления и вбейте в каждое по яйцу.",
        "Накройте крышкой и готовьте 4–6 минут, чтобы белки схватились, а желтки остались жидкими.",
        "Посолите и посыпьте свежей петрушкой перед подачей."
      ]
    },
    nutrition: { calories: 280, protein: 14, fat: 18, carbs: 12 },
    tips: {
      en: [
        "Use a lid to trap steam — it sets the egg whites without flipping.",
        "Serve straight from the pan with crusty bread to scoop up the sauce."
      ],
      ru: [
        "Используйте крышку, чтобы задержать пар — это схватит белки без переворачивания.",
        "Подавайте прямо из сковороды с хрустящим хлебом, чтобы макать в соус."
      ]
    },
    tags: ["eggs", "breakfast", "mediterranean", "tomato", "vegetarian"],
    datePublished: "2024-07-02"
  },
  {
    id: "9",
    slug: "beef-tacos",
    title: { en: "Beef Tacos", ru: "Говяжьи тако" },
    description: {
      en: "Spiced ground beef in crisp corn shells with fresh salsa, lettuce, and a squeeze of lime.",
      ru: "Острый мясной фарш в хрустящих кукурузных лепёшках со свежей сальсой, салатом и долькой лайма."
    },
    image:
      "https://images.unsplash.com/photo-1551504734-5ee1c4a1479b?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Beef tacos with salsa and lime on a plate",
      ru: "Говяжьи тако с сальсой и лаймом на тарелке"
    },
    cuisine: "mexican",
    category: "dinner",
    mealType: ["dinner"],
    diet: [],
    difficulty: "easy",
    prepTimeMinutes: 15,
    cookTimeMinutes: 15,
    totalTimeMinutes: 30,
    servings: 4,
    rating: 4.5,
    ingredients: {
      en: [
        "500 g ground beef",
        "1 onion, diced",
        "2 garlic cloves, minced",
        "1 tbsp chili powder",
        "1 tsp cumin",
        "8 corn taco shells",
        "1 tomato, diced",
        "1/4 iceberg lettuce, shredded",
        "Lime wedges and fresh cilantro"
      ],
      ru: [
        "500 г говяжьего фарша",
        "1 луковица, нарезать кубиками",
        "2 зубчика чеснока, измельчить",
        "1 ст. л. порошка чили",
        "1 ч. л. кумина",
        "8 кукурузных тако-лепёшек",
        "1 томат, нарезать",
        "1/4 кочана салата-айсберг, нашинковать",
        "Дольки лайма и свежая кинза"
      ]
    },
    steps: {
      en: [
        "Brown the ground beef with the onion in a dry pan, breaking it up as it cooks.",
        "Add the garlic, chili powder, cumin, and a splash of water and simmer for 5 minutes.",
        "Warm the taco shells in a hot oven for 2 minutes.",
        "Pile the beef into the shells and top with tomato and lettuce.",
        "Finish with lime juice and fresh cilantro and serve at once."
      ],
      ru: [
        "Обжарьте фарш с луком в сухой сковороде, разламывая комки по мере готовки.",
        "Добавьте чеснок, порошок чили, кумин и немного воды, тушите 5 минут.",
        "Разогрейте тако-лепёшки в горячей духовке 2 минуты.",
        "Наполните лепёшки фаршем и сверху выложите томат и салат.",
        "Сбрызните соком лайма и кинзой и подавайте сразу."
      ]
    },
    nutrition: { calories: 380, protein: 20, fat: 18, carbs: 30 },
    tips: {
      en: [
        "Warm the shells briefly so they stay crisp without cracking.",
        "Keep the beef a little saucy — it keeps the filling moist inside the shell."
      ],
      ru: [
        "Коротко прогрейте лепёшки, чтобы они оставались хрустящими, не ломаясь.",
        "Оставьте фарш немного сочным — так начинка будет влажной внутри лепёшки."
      ]
    },
    tags: ["tacos", "mexican", "beef", "dinner"],
    datePublished: "2024-07-22"
  },
  {
    id: "10",
    slug: "french-onion-soup",
    title: { en: "French Onion Soup", ru: "Французский луковый суп" },
    description: {
      en: "Slowly caramelised onions in a rich beef broth, topped with a toasted crouton and melted Gruyère.",
      ru: "Медленно карамелизованный лук в насыщенном говяжьем бульоне с поджаренным крутоном и расплавленным грюйером."
    },
    image:
      "https://images.unsplash.com/photo-1547592166-23ac45744acd?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "French onion soup with melted cheese on top",
      ru: "Французский луковый суп с расплавленным сыром сверху"
    },
    cuisine: "french",
    category: "soups",
    mealType: ["lunch", "dinner"],
    diet: [],
    difficulty: "medium",
    prepTimeMinutes: 15,
    cookTimeMinutes: 50,
    totalTimeMinutes: 65,
    servings: 4,
    rating: 4.6,
    ingredients: {
      en: [
        "4 large yellow onions, sliced",
        "50 g butter",
        "1 tbsp olive oil",
        "1 L beef stock",
        "150 ml white wine",
        "1 sprig thyme",
        "4 slices baguette",
        "150 g Gruyère, grated",
        "Salt and pepper"
      ],
      ru: [
        "4 крупные жёлтые луковицы, нарезать",
        "50 г сливочного масла",
        "1 ст. л. оливкового масла",
        "1 л говяжьего бульона",
        "150 мл белого вина",
        "1 веточка тимьяна",
        "4 ломтика багета",
        "150 г грюйера, натереть",
        "Соль и перец"
      ]
    },
    steps: {
      en: [
        "Melt the butter with the olive oil in a heavy pot and add the sliced onions.",
        "Cook over low heat for 40 minutes, stirring often, until deeply caramelised.",
        "Pour in the wine and let it reduce by half, then add the beef stock and thyme.",
        "Simmer for 15 minutes and season carefully with salt and pepper.",
        "Ladle the soup into bowls, top with a slice of baguette and Gruyère, and grill until bubbly."
      ],
      ru: [
        "Растопите масло с оливковым в толстостенной кастрюле и добавьте нарезанный лук.",
        "Готовьте на слабом огне 40 минут, часто помешивая, до тёмной карамелизации.",
        "Влейте вино и дайте упариться наполовину, затем добавьте бульон и тимьян.",
        "Варите 15 минут и тщательно посолите и поперчите.",
        "Разлейте суп по тарелкам, сверху положите ломтик багета и грюйер, зажарьте до пузырей."
      ]
    },
    nutrition: { calories: 230, protein: 8, fat: 9, carbs: 28 },
    tips: {
      en: [
        "Patience is everything — the onions should be jammy and brown, not burnt.",
        "Use a good beef stock; it is the backbone of the whole soup."
      ],
      ru: [
        "Терпение — главное: лук должен быть тягучим и коричневым, но не сгоревшим.",
        "Используйте хороший говяжий бульон — он основа всего супа."
      ]
    },
    tags: ["soup", "french", "onion", "comfort-food", "classic"],
    datePublished: "2024-08-10"
  },
  {
    id: "11",
    slug: "fluffy-pancakes",
    title: { en: "Fluffy Pancakes", ru: "Пышные оладьи" },
    description: {
      en: "Tall, soft American-style pancakes with a hint of vanilla, served with maple syrup and butter.",
      ru: "Высокие мягкие американские оладьи с ванилью, подаются с кленовым сиропом и маслом."
    },
    image:
      "https://images.unsplash.com/photo-1567620905732-2d1ec7ab7445?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Stack of fluffy pancakes with maple syrup",
      ru: "Стопка пышных оладий с кленовым сиропом"
    },
    cuisine: "american",
    category: "breakfast",
    mealType: ["breakfast"],
    diet: ["vegetarian"],
    difficulty: "easy",
    prepTimeMinutes: 10,
    cookTimeMinutes: 10,
    totalTimeMinutes: 20,
    servings: 4,
    rating: 4.7,
    ingredients: {
      en: [
        "200 g all-purpose flour",
        "1 tbsp sugar",
        "1 tsp baking powder",
        "1/2 tsp baking soda",
        "1/2 tsp salt",
        "250 ml buttermilk",
        "1 egg",
        "50 g melted butter",
        "1 tsp vanilla extract",
        "Maple syrup to serve"
      ],
      ru: [
        "200 г пшеничной муки",
        "1 ст. л. сахара",
        "1 ч. л. пекарского порошка",
        "1/2 ч. л. соды",
        "1/2 ч. л. соли",
        "250 мл пахты",
        "1 яйцо",
        "50 г растопленного масла",
        "1 ч. л. ванильного экстракта",
        "Кленовый сироп для подачи"
      ]
    },
    steps: {
      en: [
        "Whisk the flour, sugar, baking powder, baking soda, and salt in a bowl.",
        "In a separate bowl, whisk the buttermilk, egg, melted butter, and vanilla.",
        "Fold the wet into the dry until just combined — a few lumps are fine.",
        "Heat a griddle over medium heat and lightly grease it.",
        "Ladle the batter into rounds and flip when bubbles form on the surface, about 2 minutes per side.",
        "Serve warm with maple syrup and extra butter."
      ],
      ru: [
        "В миске смешайте муку, сахар, пекарский порошок, соду и соль.",
        "В отдельной миске взбейте пахту, яйцо, растопленное масло и ваниль.",
        "Соедините сухую и влажную смеси, перемешав до объединения — комочки допустимы.",
        "Разогрейте сковороду на среднем огне и смажьте её.",
        "Вылейте тесто кружочками и переворачивайте, когда сверху появятся пузырьки, примерно 2 минуты с каждой стороны.",
        "Подавайте тёплыми с кленовым сиропом и добавочным маслом."
      ]
    },
    nutrition: { calories: 320, protein: 8, fat: 14, carbs: 42 },
    tips: {
      en: [
        "Let the batter rest for 5 minutes for the thickest, fluffiest pancakes.",
        "Keep the heat medium — too high browns the outside before the middle cooks."
      ],
      ru: [
        "Дайте тесту постоять 5 минут — оладьи будут толстыми и пышными.",
        "Держите средний огонь — на сильном снаружи сгорит, а середина останется сырой."
      ]
    },
    tags: ["pancakes", "breakfast", "american", "sweet", "vegetarian"],
    datePublished: "2024-08-25"
  },
  {
    id: "12",
    slug: "apple-pie",
    title: { en: "Apple Pie", ru: "Яблочный пирог" },
    description: {
      en: "A classic American double-crust apple pie with cinnamon-spiced filling and a flaky butter crust.",
      ru: "Классический американский яблочный пирог с двойной хрустящей коркой и начинкой с корицей."
    },
    image:
      "https://images.unsplash.com/photo-1568571780765-9276ac8b75a2?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Apple pie with a lattice crust on a cooling rack",
      ru: "Яблочный пирог с решётчатой коркой на охладительной решётке"
    },
    cuisine: "american",
    category: "desserts",
    mealType: ["dessert"],
    diet: ["vegetarian"],
    difficulty: "medium",
    prepTimeMinutes: 25,
    cookTimeMinutes: 45,
    totalTimeMinutes: 70,
    servings: 8,
    rating: 4.8,
    ingredients: {
      en: [
        "300 g all-purpose flour",
        "200 g cold butter, cubed",
        "1 tsp salt",
        "60 ml ice water",
        "6 apples, peeled and sliced",
        "100 g sugar",
        "1 tsp cinnamon",
        "1 tbsp lemon juice"
      ],
      ru: [
        "300 г пшеничной муки",
        "200 г холодного масла, нарезать кубиками",
        "1 ч. л. соли",
        "60 мл ледяной воды",
        "6 яблок, очищенных и нарезанных",
        "100 г сахара",
        "1 ч. л. корицы",
        "1 ст. л. лимонного сока"
      ]
    },
    steps: {
      en: [
        "Rub the butter into the flour and salt until the mixture looks like breadcrumbs.",
        "Add the ice water and bring the dough together, then wrap and chill for 30 minutes.",
        "Roll out two-thirds of the dough and line a pie dish.",
        "Toss the apples with sugar, cinnamon, and lemon juice and pile into the crust.",
        "Roll the remaining dough, lay over the filling, and crimp the edges.",
        "Bake at 200 °C for 45 minutes until golden and bubbling."
      ],
      ru: [
        "Втирайте масло в муку с солью, пока смесь не станет похожа на хлебные крошки.",
        "Добавьте ледяную воду и соберите тесто, затем заверните и охладите 30 минут.",
        "Раскатайте две трети теста и выложите им форму для пирога.",
        "Перемешайте яблоки с сахаром, корицей и лимонным соком и наполните корж.",
        "Раскатайте оставшееся тесто, накройте начинку и зажмите края.",
        "Выпекайте при 200 °C 45 минут до золотистости и пузырей."
      ]
    },
    nutrition: { calories: 290, protein: 4, fat: 14, carbs: 40 },
    tips: {
      en: [
        "Keep the butter cold — it melts in the oven and creates flake.",
        "Let the pie cool for an hour before slicing so the filling sets."
      ],
      ru: [
        "Держите масло холодным — оно тает в духовке, создавая хрустящие слои.",
        "Дайте пирогу остыть час перед нарезкой, чтобы начинка схватилась."
      ]
    },
    tags: ["pie", "dessert", "apple", "american", "baking"],
    datePublished: "2024-09-12"
  },
  {
    id: "13",
    slug: "greek-salad",
    title: { en: "Greek Salad", ru: "Греческий салат" },
    description: {
      en: "Sun-ripe tomatoes, cucumber, olives, and feta dressed with oregano and the best olive oil you have.",
      ru: "Солнечные томаты, огурец, оливки и фета, заправленные орегано и лучшим оливковым маслом."
    },
    image:
      "https://images.unsplash.com/photo-1540420773420-3366772f4999?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Greek salad with feta and olives in a bowl",
      ru: "Греческий салат с фетой и оливками в миске"
    },
    cuisine: "mediterranean",
    category: "salads",
    mealType: ["lunch"],
    diet: ["vegetarian"],
    difficulty: "easy",
    prepTimeMinutes: 15,
    cookTimeMinutes: 0,
    totalTimeMinutes: 15,
    servings: 4,
    rating: 4.5,
    ingredients: {
      en: [
        "4 tomatoes, wedged",
        "1 cucumber, sliced",
        "1/2 red onion, sliced",
        "100 g Kalamata olives",
        "200 g feta block",
        "2 tbsp olive oil",
        "1 tsp dried oregano",
        "Salt and pepper"
      ],
      ru: [
        "4 томата, дольками",
        "1 огурец, нарезать",
        "1/2 красной луковицы, нарезать",
        "100 г оливок каламата",
        "200 г феты куском",
        "2 ст. л. оливкового масла",
        "1 ч. л. сушёного орегано",
        "Соль и перец"
      ]
    },
    steps: {
      en: [
        "Combine the tomatoes, cucumber, onion, and olives in a wide bowl.",
        "Season with salt, pepper, and half the oregano.",
        "Place the feta on top in a single block and sprinkle with the rest of the oregano.",
        "Drizzle with the olive oil and serve immediately so the vegetables stay crisp."
      ],
      ru: [
        "Соедините томаты, огурец, лук и оливки в широкой миске.",
        "Посолите, поперчите и посыпьте половиной орегано.",
        "Положите фету сверху целым куском и посыпьте оставшимся орегано.",
        "Сбрызните оливковым маслом и подавайте сразу, чтобы овощи оставались хрустящими."
      ]
    },
    nutrition: { calories: 220, protein: 7, fat: 16, carbs: 14 },
    tips: {
      en: [
        "Keep the feta in a block — it looks more authentic and tastes creamier.",
        "Don't cut the vegetables too small — rustic chunks are the point."
      ],
      ru: [
        "Держите фету куском — это выглядит аутентичнее и вкус насыщеннее.",
        "Не нарезайте овощи слишком мелко — крупные куски — в этом суть."
      ]
    },
    tags: ["salad", "mediterranean", "greek", "vegetarian", "fresh"],
    datePublished: "2024-09-30"
  },
  {
    id: "14",
    slug: "miso-soup",
    title: { en: "Miso Soup", ru: "Суп мисо" },
    description: {
      en: "A delicate Japanese soup with a clear dashi base, tender tofu, wakame, and a spoonful of white miso.",
      ru: "Деликатный японский суп на прозрачном бульоне даси с нежным тофу, вакамэ и ложкой белого мисо."
    },
    image:
      "https://images.unsplash.com/photo-1547592166-23ac45744acd?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Bowl of miso soup with tofu and seaweed",
      ru: "Тарелка супа мисо с тофу и водорослями"
    },
    cuisine: "japanese",
    category: "soups",
    mealType: ["lunch", "dinner"],
    diet: ["vegan", "dairy-free"],
    difficulty: "easy",
    prepTimeMinutes: 10,
    cookTimeMinutes: 10,
    totalTimeMinutes: 20,
    servings: 4,
    rating: 4.3,
    ingredients: {
      en: [
        "1 L dashi stock",
        "3 tbsp white miso paste",
        "150 g soft tofu, cubed",
        "2 tbsp dried wakame",
        "2 spring onions, sliced"
      ],
      ru: [
        "1 л бульона даси",
        "3 ст. л. пасты белого мисо",
        "150 г мягкого тофу, кубиками",
        "2 ст. л. сушёного вакамэ",
        "2 стебля зелёного лука, нарезать"
      ]
    },
    steps: {
      en: [
        "Warm the dashi in a pot without bringing it to a boil.",
        "Soak the wakame in cold water for 5 minutes until expanded.",
        "Whisk the miso paste with a ladle of warm dashi, then stir back into the pot.",
        "Add the tofu and drained wakame and heat gently for 2 minutes.",
        "Ladle into bowls and top with sliced spring onion."
      ],
      ru: [
        "Подогрейте даси в кастрюле, не доводя до кипения.",
        "Замочите вакамэ в холодной воде на 5 минут до набухания.",
        "Взбейте пасту мисо с половником тёплого даси, затем верните в кастрюлю.",
        "Добавьте тофу и отжатый вакамэ, прогрейте 2 минуты.",
        "Разлейте по тарелкам и посыпьте нарезанным зелёным луком."
      ]
    },
    nutrition: { calories: 110, protein: 6, fat: 3, carbs: 12 },
    tips: {
      en: [
        "Never boil the miso — it kills the flavour and the good bacteria.",
        "Use white miso for a mellow taste, or mix in a little red miso for depth."
      ],
      ru: [
        "Не кипятите мисо — это убьёт вкус и полезные бактерии.",
        "Используйте белое мисо для мягкого вкуса или добавьте немного красного для насыщенности."
      ]
    },
    tags: ["soup", "japanese", "miso", "vegan", "light"],
    datePublished: "2024-10-15"
  },
  {
    id: "15",
    slug: "mushroom-risotto",
    title: { en: "Mushroom Risotto", ru: "Грибное ризотто" },
    description: {
      en: "Creamy arborio risotto with pan-seared mushrooms, white wine, and a final fold of butter and Parmesan.",
      ru: "Кремовое ризотто из арборио с обжаренными грибами, белым вином и финальным вмешиванием масла и пармезана."
    },
    image:
      "https://images.unsplash.com/photo-1476124369491-e7addf5db371?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Mushroom risotto in a shallow bowl with parmesan",
      ru: "Грибное ризотто в неглубокой тарелке с пармезаном"
    },
    cuisine: "italian",
    category: "dinner",
    mealType: ["dinner"],
    diet: ["vegetarian"],
    difficulty: "medium",
    prepTimeMinutes: 10,
    cookTimeMinutes: 25,
    totalTimeMinutes: 35,
    servings: 4,
    rating: 4.7,
    ingredients: {
      en: [
        "300 g arborio rice",
        "300 g mixed mushrooms, sliced",
        "1 onion, diced",
        "2 garlic cloves, minced",
        "150 ml white wine",
        "1 L vegetable stock, hot",
        "50 g butter",
        "50 g Parmesan, grated",
        "1 tbsp olive oil",
        "Fresh thyme"
      ],
      ru: [
        "300 г риса арборио",
        "300 г грибов, нарезать",
        "1 луковица, нарезать кубиками",
        "2 зубчика чеснока, измельчить",
        "150 мл белого вина",
        "1 л овощного бульона, горячего",
        "50 г сливочного масла",
        "50 г пармезана, натереть",
        "1 ст. л. оливкового масла",
        "Свежий тимьян"
      ]
    },
    steps: {
      en: [
        "Sear the mushrooms in olive oil until golden and set aside.",
        "Soften the onion and garlic in the same pan, then add the rice and toast for a minute.",
        "Pour in the wine and stir until fully absorbed.",
        "Add the stock a ladle at a time, stirring until each addition is absorbed.",
        "After 18 minutes, fold in the mushrooms, butter, and Parmesan.",
        "Cover, rest for 2 minutes, and serve topped with fresh thyme."
      ],
      ru: [
        "Обжарьте грибы на оливковом масле до золотистости и отложите.",
        "В той же сковороде смягчите лук и чеснок, добавьте рис и прогрейте минуту.",
        "Влейте вино и перемешивайте до полного впитывания.",
        "Добавляйте бульон по половнику, перемешивая, пока каждый не впитается.",
        "Через 18 минут вмешайте грибы, масло и пармезан.",
        "Накройте, дайте постоять 2 минуты и подавайте, посыпав тимьяном."
      ]
    },
    nutrition: { calories: 420, protein: 10, fat: 16, carbs: 56 },
    tips: {
      en: [
        "Stir often — it releases the rice's starch for a creamy finish.",
        "Taste at the end — the rice should be al dente, not soft."
      ],
      ru: [
        "Часто помешивайте — это выделяет крахмал из риса для кремовой текстуры.",
        "Попробуйте в конце — рис должен быть аль денте, не разваренным."
      ]
    },
    tags: ["risotto", "italian", "mushrooms", "vegetarian", "creamy"],
    datePublished: "2024-11-05"
  },
  {
    id: "16",
    slug: "latvian-rye-bread-dessert",
    title: { en: "Latvian Rye Bread Dessert", ru: "Латышский десерт из ржаного хлеба" },
    description: {
      en: "A layered Latvian classic with caramelised rye bread crumbs, whipped cream, and lingonberry jam.",
      ru: "Многослойный латышский десерт с карамелизованными крошками ржаного хлеба, сливками и брусничным вареньем."
    },
    image:
      "https://images.unsplash.com/photo-1488477181946-6428a0291777?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Layered rye bread dessert with cream and berries in a glass",
      ru: "Многослойный десерт из ржаного хлеба с сливками и ягодами в стакане"
    },
    cuisine: "latvian",
    category: "desserts",
    mealType: ["dessert"],
    diet: ["vegetarian"],
    difficulty: "easy",
    prepTimeMinutes: 15,
    cookTimeMinutes: 15,
    totalTimeMinutes: 30,
    servings: 4,
    rating: 4.4,
    ingredients: {
      en: [
        "150 g rye bread, crusts removed",
        "50 g butter",
        "3 tbsp sugar",
        "200 ml double cream",
        "1 tsp vanilla extract",
        "4 tbsp lingonberry or cranberry jam"
      ],
      ru: [
        "150 г ржаного хлеба, без корок",
        "50 г сливочного масла",
        "3 ст. л. сахара",
        "200 мл жирных сливок",
        "1 ч. л. ванильного экстракта",
        "4 ст. л. брусничного или клюквенного варенья"
      ]
    },
    steps: {
      en: [
        "Tear the rye bread into fine crumbs in a food processor.",
        "Melt the butter in a pan and toast the crumbs with the sugar until caramelised.",
        "Whip the cream with the vanilla until soft peaks form.",
        "Layer the crumbs, whipped cream, and jam in glasses.",
        "Chill for 20 minutes before serving."
      ],
      ru: [
        "Измельчите ржаной хлеб в мелкие крошки в процессоре.",
        "Растопите масло на сковороде и обжарьте крошки с сахаром карамелизации.",
        "Взбейте сливки с ванилью до мягких пиков.",
        "Выложите крошки, сливки и варенье в стаканах слоями.",
        "Охладите 20 минут перед подачей."
      ]
    },
    nutrition: { calories: 280, protein: 6, fat: 12, carbs: 38 },
    tips: {
      en: [
        "Use a dense, dark rye — the flavour is closer to the original.",
        "Whip the cream to soft peaks only — it should still flow a little."
      ],
      ru: [
        "Используйте плотный тёмный ржаной хлеб — вкус будет ближе к оригинату.",
        "Взбивайте сливки только до мягких пиков — они должны ещё чуть теч."
      ]
    },
    tags: ["dessert", "latvian", "rye-bread", "cranberry", "vegetarian"],
    datePublished: "2024-11-22"
  },
  {
    id: "17",
    slug: "falafel-bowl",
    title: { en: "Falafel Bowl", ru: "Боул с фалафелем" },
    description: {
      en: "Crispy baked falafel over fluffy quinoa with cucumber, tomato, and a drizzle of tahini-lemon dressing.",
      ru: "Хрустящий запечённый фалафель на рыхлом киноа с огурцом, томатом и тахини-лимонной заправкой."
    },
    image:
      "https://images.unsplash.com/photo-1593001874117-c99c800e3eb7?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Falafel bowl with quinoa and vegetables in a bowl",
      ru: "Боул с фалафелем, киноа и овощами в миске"
    },
    cuisine: "mediterranean",
    category: "vegetarian",
    mealType: ["lunch", "dinner"],
    diet: ["vegan", "gluten-free"],
    difficulty: "easy",
    prepTimeMinutes: 15,
    cookTimeMinutes: 20,
    totalTimeMinutes: 35,
    servings: 4,
    rating: 4.6,
    ingredients: {
      en: [
        "200 g dried chickpeas, soaked overnight",
        "1/2 onion, diced",
        "2 garlic cloves",
        "1 tsp cumin",
        "1 tsp coriander",
        "2 tbsp olive oil",
        "200 g quinoa",
        "1 cucumber, diced",
        "2 tbsp tahini",
        "1 tbsp lemon juice"
      ],
      ru: [
        "200 г сушёного нута, замочить на ночь",
        "1/2 луковицы, нарезать",
        "2 зубчика чеснока",
        "1 ч. л. кумина",
        "1 ч. л. кориандра",
        "2 ст. л. оливкового масла",
        "200 г киноа",
        "1 огурец, нарезать",
        "2 ст. л. тахини",
        "1 ст. л. лимонного сока"
      ]
    },
    steps: {
      en: [
        "Blend the chickpeas, onion, garlic, cumin, and coriander into a rough paste.",
        "Shape into 16 small patties and arrange on a tray.",
        "Bake at 200 °C for 20 minutes, flipping halfway, until crisp.",
        "Cook the quinoa according to the package directions and fluff with a fork.",
        "Whisk the tahini with lemon juice and a splash of water to loosen.",
        "Divide the quinoa, falafel, and cucumber into bowls and drizzle with the dressing."
      ],
      ru: [
        "Измельчите нут, лук, чеснок, кумин и кориандр в грубую пасту.",
        "Сформируйте 16 маленьких котлет и выложите на противень.",
        "Запекайте при 200 °C 20 минут, переворачивая в середине, до хруста.",
        "Отварите киноа по инструкции на упаковке и взрыхлите вилкой.",
        "Взбейте тахини с лимонным соком и водой до лёгкой консистенции.",
        "Разложите по мискам киноа, фалафель и огурец, сбрызните соусом."
      ]
    },
    nutrition: { calories: 380, protein: 14, fat: 16, carbs: 44 },
    tips: {
      en: [
        "Don't use canned chickpeas — soaked dried ones give the right texture.",
        "Baking instead of frying keeps the falafel crisp without too much oil."
      ],
      ru: [
        "Не используйте консервированный нут — сушёный даёт правильную текстуру.",
        "Запекание вместо жарки сохраняет хруст без избытка масла."
      ]
    },
    tags: ["falafel", "mediterranean", "vegan", "bowl", "healthy"],
    datePublished: "2024-12-10"
  },
  {
    id: "18",
    slug: "salmon-teriyaki",
    title: { en: "Salmon Teriyaki", ru: "Лосось терияки" },
    description: {
      en: "Glazed salmon fillets in a glossy soy-mirin sauce, served over steamed rice with sesame and spring onion.",
      ru: "Лосось в блестящем соево-мириновом соусе, подаётся на пару с рисом, кунжутом и зелёным луком."
    },
    image:
      "https://images.unsplash.com/photo-1467003909585-2f8a72700288?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Teriyaki salmon fillet with sesame and spring onion",
      ru: "Филе лосося в терияки с кунжутом и зелёным луком"
    },
    cuisine: "japanese",
    category: "dinner",
    mealType: ["dinner"],
    diet: ["pescatarian", "dairy-free"],
    difficulty: "easy",
    prepTimeMinutes: 10,
    cookTimeMinutes: 15,
    totalTimeMinutes: 25,
    servings: 2,
    rating: 4.7,
    ingredients: {
      en: [
        "2 salmon fillets, skin on",
        "3 tbsp soy sauce",
        "2 tbsp mirin",
        "1 tbsp sugar",
        "1 garlic clove, grated",
        "1 tsp grated ginger",
        "200 g short-grain rice",
        "1 tbsp sesame seeds",
        "2 spring onions, sliced"
      ],
      ru: [
        "2 филе лосося, с кожей",
        "3 ст. л. соевого соуса",
        "2 ст. л. мирина",
        "1 ст. л. сахара",
        "1 зубчик чеснока, натереть",
        "1 ч. л. натёртого имбиря",
        "200 г круглозерного риса",
        "1 ст. л. кунжута",
        "2 стебля зелёного лука, нарезать"
      ]
    },
    steps: {
      en: [
        "Combine the soy sauce, mirin, sugar, garlic, and ginger in a small bowl.",
        "Sear the salmon skin-side down in a hot pan for 3 minutes until crisp.",
        "Flip the fillets, pour in the sauce, and spoon it over for 4 minutes until glazed.",
        "Cook the rice according to the package directions.",
        "Serve the salmon over rice, topped with sesame and spring onion."
      ],
      ru: [
        "Соедините соевый соус, мирин, сахар, чеснок и имбирь в маленькой миске.",
        "Обжарьте лосось кожей вниз на горячей сковороде 3 минуты до хруста.",
        "Переверните филе, влейте соус и поливайте им рыбу 4 минуты до глазури.",
        "Отварите рис по инструкции на упаковке.",
        "Подавайте лосося на рисе, посыпав кунжутом и зелёным луком."
      ]
    },
    nutrition: { calories: 420, protein: 32, fat: 22, carbs: 12 },
    tips: {
      en: [
        "Start the salmon in a cold pan skin-down for an even crisp skin.",
        "Don't overcook — the centre should still be just translucent."
      ],
      ru: [
        "Начните лосось в холодной сковороде кожей вниз для равномерного хруста.",
        "Не переваривайте — середина должна быть чуть полупрозрачной."
      ]
    },
    tags: ["salmon", "japanese", "teriyaki", "pescatarian", "dinner"],
    datePublished: "2024-12-28"
  },
  {
    id: "19",
    slug: "lentil-soup",
    title: { en: "Lentil Soup", ru: "Чечевичный суп" },
    description: {
      en: "A hearty, low-fat vegan soup with red lentils, carrot, and a hint of cumin — ready in under an hour.",
      ru: "Сытный обезжиренный веганский суп с красной чечевицей, морковью и ноткой кумина — готовится меньше часа."
    },
    image:
      "https://images.unsplash.com/photo-1547592166-23ac45744acd?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Bowl of lentil soup with a swirl of olive oil",
      ru: "Тарелка чечевичного супа с завитком оливкового масла"
    },
    cuisine: "mediterranean",
    category: "soups",
    mealType: ["lunch", "dinner"],
    diet: ["vegan", "gluten-free", "dairy-free"],
    difficulty: "easy",
    prepTimeMinutes: 10,
    cookTimeMinutes: 30,
    totalTimeMinutes: 40,
    servings: 6,
    rating: 4.5,
    ingredients: {
      en: [
        "200 g red lentils, rinsed",
        "1 onion, diced",
        "2 carrots, diced",
        "2 garlic cloves, minced",
        "1 tsp cumin",
        "1 L vegetable stock",
        "1 tbsp lemon juice",
        "Olive oil to serve"
      ],
      ru: [
        "200 г красной чечевицы, промыть",
        "1 луковица, нарезать",
        "2 моркови, нарезать",
        "2 зубчика чеснока, измельчить",
        "1 ч. л. кумина",
        "1 л овощного бульона",
        "1 ст. л. лимонного сока",
        "Оливковое масло для подачи"
      ]
    },
    steps: {
      en: [
        "Soften the onion and carrot in a little oil for 5 minutes.",
        "Add the garlic and cumin and stir for a minute.",
        "Stir in the lentils and the vegetable stock and bring to a boil.",
        "Simmer for 25 minutes until the lentils break down completely.",
        "Blend half the soup for a creamy texture, stir in the lemon juice, and season."
      ],
      ru: [
        "Смягчите лук и морковь на небольшом количестве масла 5 минут.",
        "Добавьте чеснок и кумин, перемешивайте минуту.",
        "Добавьте чечевицу и овощной бульон, доведите до кипения.",
        "Варите 25 минут, пока чечевица полностью не разварится.",
        "Пюррируйте половину супа для кремовой текстуры, влейте лимонный сок и посолите."
      ]
    },
    nutrition: { calories: 230, protein: 12, fat: 5, carbs: 32 },
    tips: {
      en: [
        "A splash of lemon at the end lifts the lentils — don't skip it.",
        "Rinse the lentils well — it removes the starch that can scorch."
      ],
      ru: [
        "Ложка лимона в конце оживляет чечевицу — не пропускайте этот шаг.",
        "Тщательно промойте чечевицу — это удаляет крахмал, который может подгореть."
      ]
    },
    tags: ["soup", "lentils", "mediterranean", "vegan", "healthy"],
    datePublished: "2025-01-15"
  },
  {
    id: "20",
    slug: "tiramisu",
    title: { en: "Tiramisu", ru: "Тирамису" },
    description: {
      en: "Layers of espresso-soaked savoiardi, mascarpone cream, and a dusting of cocoa — the definitive Italian dessert.",
      ru: "Многослойный десерт из печенья савоярди, вымоченного в эспрессо, крема из маскарпоне и просыпки какао — классический итальянский десерт."
    },
    image:
      "https://images.unsplash.com/photo-1571877227200-a0d98ea607e9?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Slice of tiramisu dusted with cocoa on a plate",
      ru: "Кусок тирамису, посыпанный какао, на тарелке"
    },
    cuisine: "italian",
    category: "desserts",
    mealType: ["dessert"],
    diet: ["vegetarian"],
    difficulty: "medium",
    prepTimeMinutes: 20,
    cookTimeMinutes: 0,
    totalTimeMinutes: 20,
    servings: 8,
    rating: 4.9,
    ingredients: {
      en: [
        "300 g savoiardi biscuits",
        "300 ml strong espresso, cooled",
        "250 g mascarpone, room temperature",
        "3 egg yolks",
        "80 g sugar",
        "300 ml double cream",
        "2 tbsp cocoa powder",
        "2 tbsp marsala or dark rum (optional)"
      ],
      ru: [
        "300 г печенья савоярди",
        "300 мл крепкого эспрессо, охлаждённого",
        "250 г маскарпоне, комнатной температуры",
        "3 яичных желтка",
        "80 г сахара",
        "300 мл жирных сливок",
        "2 ст. л. какао-порошка",
        "2 ст. л. марсалы или тёмного рума (по желанию)"
      ]
    },
    steps: {
      en: [
        "Whisk the yolks with the sugar until pale and doubled.",
        "Fold in the mascarpone until smooth, then the double cream until the mixture holds soft peaks.",
        "Dip each biscuit in the espresso for a second and arrange in a dish.",
        "Spread half the mascarpone cream over the biscuits and repeat with another layer.",
        "Chill for at least 4 hours, then dust with cocoa before serving."
      ],
      ru: [
        "Взбейте желтки с сахаром до светлой и удвоенной массы.",
        "Вмешайте маскарпоне до гладкости, затем сливки до мягких пиков.",
        "Окуните каждое печенье в эспрессо на секунду и выложите в форму.",
        "Размажьте половину крема поверх печенья и повторите слоем.",
        "Охладите не менее 4 часов, затем присыпьте какао перед подачей."
      ]
    },
    nutrition: { calories: 320, protein: 6, fat: 18, carbs: 36 },
    tips: {
      en: [
        "Don't soak the biscuits long — they should be moist, not soggy.",
        "Chill overnight if you can — the flavours deepen beautifully."
      ],
      ru: [
        "Не мочите печенье долго — оно должно быть влажным, а не мокрым.",
        "Охладите на ночь, если можно — вкусы глубже раскроются."
      ]
    },
    tags: ["dessert", "italian", "tiramisu", "coffee", "vegetarian"],
    datePublished: "2025-02-04"
  },
  {
    id: "21",
    slug: "ratatouille",
    title: { en: "Ratatouille", ru: "Рататуй" },
    description: {
      en: "A sun-drenched Provençal stew of eggplant, zucchini, peppers, and tomato simmered with herbes de Provence.",
      ru: "Пронизанный солнцем провансальский рагу из баклажана, цуккини, перцев и томатов, тушёный с прованскими травами."
    },
    image:
      "https://images.unsplash.com/photo-1576444356170-66073046b1bc?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Ratatouille with eggplant and zucchini in a baking dish",
      ru: "Рататуй с баклажаном и цуккини в форме для запекания"
    },
    cuisine: "french",
    category: "vegetarian",
    mealType: ["dinner"],
    diet: ["vegan", "gluten-free", "low-carb"],
    difficulty: "medium",
    prepTimeMinutes: 20,
    cookTimeMinutes: 45,
    totalTimeMinutes: 65,
    servings: 6,
    rating: 4.4,
    ingredients: {
      en: [
        "1 eggplant, cubed",
        "1 zucchini, sliced",
        "1 red bell pepper, sliced",
        "1 onion, diced",
        "3 garlic cloves, minced",
        "400 g chopped tomatoes",
        "2 tbsp olive oil",
        "1 tbsp herbes de Provence",
        "Salt and pepper"
      ],
      ru: [
        "1 баклажан, кубиками",
        "1 цуккини, нарезать",
        "1 красный болгарский перец, нарезать",
        "1 луковица, нарезать",
        "3 зубчика чеснока, измельчить",
        "400 г нарезанных томатов",
        "2 ст. л. оливкового масла",
        "1 ст. л. прованских трав",
        "Соль и перец"
      ]
    },
    steps: {
      en: [
        "Salt the eggplant cubes and rest for 15 minutes, then pat dry.",
        "Sear the eggplant and zucchini in olive oil until golden and set aside.",
        "Soften the onion and pepper, then add the garlic and tomatoes.",
        "Return the eggplant and zucchini to the pan with the herbes de Provence.",
        "Cover and simmer gently for 30 minutes until everything is soft and harmonious.",
        "Season and serve warm with crusty bread or over rice."
      ],
      ru: [
        "Посолите баклажан и дайте постоять 15 минут, затем обсушите.",
        "Обжарьте баклажан и цуккини на оливковом масле до золотистости и отложите.",
        "Смягчите лук и перец, добавьте чеснок и томаты.",
        "Верните баклажан и цуккини в сковороду с прованскими травами.",
        "Накройте и тушите 30 минут до полной мягкости и гармонии вкуса.",
        "Посолите и подавайте тёплым с хрустящим хлебом или с рисом."
      ]
    },
    nutrition: { calories: 180, protein: 5, fat: 10, carbs: 18 },
    tips: {
      en: [
        "Salting the eggplant first removes bitterness and keeps it from going soggy.",
        "Cook it gently — a hard boil breaks the vegetables into mush."
      ],
      ru: [
        "Соль на баклажане убирает горечь и не даёт ему стать клёклым.",
        "Тушите мягко — сильное кипение разварит овощи в пюре."
      ]
    },
    tags: ["ratatouille", "french", "vegetables", "vegan", "gluten-free"],
    datePublished: "2025-02-22"
  },
  {
    id: "22",
    slug: "hummus-plate",
    title: { en: "Hummus Plate", ru: "Тарелка хумуса" },
    description: {
      en: "Silky hummus swirled with olive oil and paprika, with warm flatbreads and crisp crudités for dipping.",
      ru: "Шёлковый хумус, завитый оливковым маслом и паприкой, с тёплыми лепёшками и свежими овощами для макания."
    },
    image:
      "https://images.unsplash.com/photo-1593001874117-c99c800e3eb7?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Hummus plate with olive oil, paprika and vegetables",
      ru: "Тарелка хумуса с оливковым маслом, паприкой и овощами"
    },
    cuisine: "mediterranean",
    category: "vegetarian",
    mealType: ["lunch", "snack"],
    diet: ["vegan", "gluten-free"],
    difficulty: "easy",
    prepTimeMinutes: 10,
    cookTimeMinutes: 0,
    totalTimeMinutes: 10,
    servings: 4,
    rating: 4.3,
    ingredients: {
      en: [
        "400 g cooked chickpeas, drained",
        "3 tbsp tahini",
        "2 garlic cloves",
        "1 lemon, juiced",
        "3 tbsp olive oil",
        "1/2 tsp smoked paprika",
        "Carrot and cucumber sticks to serve"
      ],
      ru: [
        "400 г отваренного нута, слить",
        "3 ст. л. тахини",
        "2 зубчика чеснока",
        "1 лимон, сок",
        "3 ст. л. оливкового масла",
        "1/2 ч. л. копчёной паприки",
        "Палочки моркови и огурца для подачи"
      ]
    },
    steps: {
      en: [
        "Blend the chickpeas, tahini, garlic, lemon juice, and 2 tbsp olive oil until silky.",
        "Scrape down and blend again, adding cold water a spoon at a time until loos.",
        "Spoon into a shallow bowl and swirl the top with the back of a spoon.",
        "Drizzle with the remaining olive oil and dust with smoked paprika.",
        "Serve with carrot and cucumber sticks."
      ],
      ru: [
        "Измельчите нут, тахини, чеснок, лимонный сок и 2 ст. л. масла до шёлковой текстуры.",
        "Соскоблите стенки и измельчите снова, добавляя холодную воду по ложке до лёгкой консистенции.",
        "Выложите в неглубокую миску и сделайте завиток ложкой сверху.",
        "Сбрызните оставшимся маслом и посыпьте паприкой.",
        "Подавайте с палочками моркови и огурца."
      ]
    },
    nutrition: { calories: 210, protein: 8, fat: 12, carbs: 20 },
    tips: {
      en: [
        "Blend the hummus long — three minutes makes it truly silky.",
        "Cold water loosens it better than oil — keep adding until smooth."
      ],
      ru: [
        "Взбивайте хумус подолгу — три минуты дают по-настоящему шёлковую текстуру.",
        "Холодная вода разжижает лучше масла — добавляйте, пока масса не станет гладкой."
      ]
    },
    tags: ["hummus", "mediterranean", "vegan", "dip", "gluten-free"],
    datePublished: "2025-03-12"
  },
  {
    id: "23",
    slug: "cheesecake",
    title: { en: "Cheesecake", ru: "Чизкейк" },
    description: {
      en: "A rich New York-style baked cheesecake on a buttery graham crust, with a vanilla bean and sour cream topping.",
      ru: "Пышный запечённый чизкейк в нью-йоркском стиле на масляной бисквитной основе с ванилью и сметанной заливкой."
    },
    image:
      "https://images.unsplash.com/photo-1533134242443-d4fd215305ad?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Slice of baked cheesecake with a graham crust",
      ru: "Кусок запечённого чизкейка с бисквитной коржой"
    },
    cuisine: "american",
    category: "desserts",
    mealType: ["dessert"],
    diet: ["vegetarian"],
    difficulty: "medium",
    prepTimeMinutes: 20,
    cookTimeMinutes: 50,
    totalTimeMinutes: 70,
    servings: 8,
    rating: 4.8,
    ingredients: {
      en: [
        "200 g graham crackers, crushed",
        "100 g butter, melted",
        "600 g cream cheese, room temperature",
        "200 g sugar",
        "3 eggs",
        "1 vanilla pod, scraped",
        "150 ml sour cream"
      ],
      ru: [
        "200 г бисквитного печенья, измельчить",
        "100 г сливочного масла, растопить",
        "600 г сливочного сыра, комнатной температуры",
        "200 г сахара",
        "3 яйца",
        "1 стручок ванили, выскрести",
        "150 мл сметаны"
      ]
    },
    steps: {
      en: [
        "Mix the crushed crackers with the melted butter and press into a springform tin.",
        "Beat the cream cheese with the sugar and vanilla until smooth.",
        "Add the eggs one at a time, beating gently after each.",
        "Pour the filling over the crust and bake at 160 °C for 45 minutes.",
        "Mix the sour cream with a spoon of the pan juices, spread on top, and bake 5 more minutes.",
        "Cool fully in the tin, then chill before serving."
      ],
      ru: [
        "Смешайте измельчённое печенье с растопленным маслом и уложите в разъёмную форму.",
        "Взбейте сливочный сыр с сахаром и ванилью до гладкости.",
        "Добавьте яйца по одному, аккуратно взбивая после каждого.",
        "Вылейте начинку на корж и выпекайте при 160 °C 45 минут.",
        "Смешайте сметану с ложкой сока из формы, намажьте сверху и пеките ещё 5 минут.",
        "Полностью охладите в форме, затем охладите перед подачей."
      ]
    },
    nutrition: { calories: 350, protein: 7, fat: 22, carbs: 32 },
    tips: {
      en: [
        "All ingredients at room temperature — lumps are the enemy of a smooth filling.",
        "Cool slowly in the oven with the door cracked to avoid cracks on top."
      ],
      ru: [
        "Все ингредиенты комнатной температуры — комки враг гладкой начинки.",
        "Охлаждайте медлено в духовке с приоткрытой дверцей, чтобы избежать трещин."
      ]
    },
    tags: ["dessert", "cheesecake", "american", "baking", "sweet"],
    datePublished: "2025-04-02"
  },
  {
    id: "24",
    slug: "berry-smoothie",
    title: { en: "Berry Smoothie", ru: "Ягодный смузи" },
    description: {
      en: "A bright, thick smoothie of frozen mixed berries, banana, and a splash of orange juice — a five-minute breakfast.",
      ru: "Освежающий густой смузи из замороженных ягод, банана и всплеска апельсинового сока — пятиминутный завтрак."
    },
    image:
      "https://images.unsplash.com/photo-1505252585461-04db1eb84625?auto=format&fit=crop&w=1200&q=80",
    imageAlt: {
      en: "Berry smoothie in a glass with fresh fruit on top",
      ru: "Ягодный смузи в стакане со свежими фруктами сверху"
    },
    cuisine: "american",
    category: "drinks",
    mealType: ["breakfast", "snack", "drink"],
    diet: ["vegetarian", "gluten-free"],
    difficulty: "easy",
    prepTimeMinutes: 5,
    cookTimeMinutes: 0,
    totalTimeMinutes: 5,
    servings: 2,
    rating: 4.6,
    ingredients: {
      en: [
        "300 g frozen mixed berries",
        "1 banana",
        "150 ml orange juice",
        "100 g greek yogurt",
        "1 tbsp honey"
      ],
      ru: [
        "300 г замороженных смешанных ягод",
        "1 банан",
        "150 мл апельсинового сока",
        "100 г греческого йогурта",
        "1 ст. л. мёда"
      ]
    },
    steps: {
      en: [
        "Pile the frozen berries and banana into a blender.",
        "Add the orange juice, yogurt, and honey.",
        "Blend until thick and smooth, about a minute.",
        "Pour into two glasses and serve at once."
      ],
      ru: [
        "Положите замороженные ягоды и банан в блендер.",
        "Добавьте апельсиновый сок, йогурт и мёд.",
        "Взбивайте до густой и гладкой массы около минуты.",
        "Разлейте по двум стаканам и подавайте сразу."
      ]
    },
    nutrition: { calories: 180, protein: 5, fat: 2, carbs: 38 },
    tips: {
      en: [
        "Frozen berries make the smoothie thick without watering it down with ice.",
        "Add a handful of spinach for a green boost — it won't change the berry flavour."
      ],
      ru: [
        "Замороженные ягоды делают смузи густым, не разжижая его льдом.",
        "Добавьте горсть шпината для полезности — на вкус ягод это не повлияет."
      ]
    },
    tags: ["smoothie", "drinks", "berries", "vegetarian", "breakfast"],
    datePublished: "2025-05-20"
  }
];
