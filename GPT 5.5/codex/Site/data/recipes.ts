import type { Locale } from "@/i18n/config";

export type LocalizedText = Record<Locale, string>;
export type Difficulty = "easy" | "medium" | "hard";
export type Diet = "vegetarian" | "vegan" | "pescatarian" | "omnivore" | "gluten-free";

export const categories = [
  "breakfast",
  "lunch",
  "dinner",
  "desserts",
  "soups",
  "salads",
  "vegetarian",
  "quick-meals",
  "baking",
  "drinks"
] as const;
export type Category = (typeof categories)[number];

export const cuisines = [
  "italian",
  "french",
  "georgian",
  "japanese",
  "mexican",
  "indian",
  "mediterranean",
  "ukrainian",
  "latvian",
  "american"
] as const;
export type Cuisine = (typeof cuisines)[number];

export type Recipe = {
  id: string;
  slug: string;
  title: LocalizedText;
  description: LocalizedText;
  image: string;
  imageAlt: LocalizedText;
  cuisine: Cuisine;
  category: Category;
  mealType: string[];
  diet: Diet[];
  difficulty: Difficulty;
  prepTimeMinutes: number;
  cookTimeMinutes: number;
  totalTimeMinutes: number;
  servings: number;
  rating: number;
  ingredients: Record<Locale, string[]>;
  steps: Record<Locale, string[]>;
  nutrition: { calories: number; protein: number; fat: number; carbs: number };
  tips: Record<Locale, string[]>;
  tags: string[];
  publishedAt: string;
};

type Seed = Omit<
  Recipe,
  "id" | "imageAlt" | "totalTimeMinutes" | "nutrition" | "publishedAt"
> & { calories: number; publishedAt?: string };

const makeRecipe = (seed: Seed, index: number): Recipe => ({
  ...seed,
  id: String(index + 1),
  imageAlt: {
    en: `Finished ${seed.title.en} on a ceramic plate`,
    ru: `Готовое блюдо «${seed.title.ru}» на керамической тарелке`
  },
  totalTimeMinutes: seed.prepTimeMinutes + seed.cookTimeMinutes,
  nutrition: {
    calories: seed.calories,
    protein: Math.round(seed.calories / 23),
    fat: Math.round(seed.calories / 38),
    carbs: Math.round(seed.calories / 13)
  },
  publishedAt: seed.publishedAt ?? `2026-0${(index % 5) + 1}-${String((index % 25) + 1).padStart(2, "0")}`
});

const photo = (id: string) =>
  `https://images.unsplash.com/${id}?auto=format&fit=crop&w=1200&q=85`;

const seeds: Seed[] = [
  {
    slug: "classic-carbonara",
    title: { en: "Classic Carbonara", ru: "Классическая карбонара" },
    description: { en: "Silky eggs, pecorino and crisp guanciale coat every strand of spaghetti.", ru: "Шёлковые яйца, пекорино и хрустящая гуанчале обволакивают каждую нить спагетти." },
    image: photo("photo-1612874742237-6526221588e3"), cuisine: "italian", category: "dinner", mealType: ["dinner"], diet: ["omnivore"], difficulty: "medium", prepTimeMinutes: 10, cookTimeMinutes: 15, servings: 2, rating: 4.9, calories: 680,
    ingredients: { en: ["200 g spaghetti", "100 g guanciale", "2 egg yolks + 1 egg", "60 g pecorino romano"], ru: ["200 г спагетти", "100 г гуанчале", "2 желтка и 1 яйцо", "60 г пекорино романо"] },
    steps: { en: ["Boil pasta in well-salted water until just al dente.", "Fry guanciale slowly until golden and reserve its fat.", "Toss hot pasta with egg-cheese mixture, fat and a splash of pasta water off the heat."], ru: ["Отварите спагетти в хорошо посоленной воде до состояния аль денте.", "Медленно обжарьте гуанчале до золотистой корочки и сохраните жир.", "Снимите пасту с огня и смешайте с яйцами, сыром, жиром и небольшим количеством воды от варки."] },
    tips: { en: ["Never add cream.", "Work off the heat to keep the sauce glossy."], ru: ["Не добавляйте сливки.", "Смешивайте вне огня, чтобы соус остался глянцевым."] }, tags: ["pasta", "italian", "паста", "италия"]
  },
  {
    slug: "margherita-pizza",
    title: { en: "Margherita Pizza", ru: "Пицца «Маргарита»" },
    description: { en: "A blistered crust, bright tomato, basil and melting mozzarella.", ru: "Пузырчатая корочка, яркий томатный соус, базилик и тающая моцарелла." },
    image: photo("photo-1513104890138-7c749659a591"), cuisine: "italian", category: "dinner", mealType: ["dinner", "lunch"], diet: ["vegetarian"], difficulty: "medium", prepTimeMinutes: 20, cookTimeMinutes: 12, servings: 2, rating: 4.8, calories: 620,
    ingredients: { en: ["1 pizza dough", "120 g crushed tomatoes", "150 g mozzarella", "Fresh basil"], ru: ["1 порция теста для пиццы", "120 г протёртых томатов", "150 г моцареллы", "Свежий базилик"] },
    steps: { en: ["Heat oven and tray as hot as possible.", "Stretch dough and spread a thin layer of tomato.", "Bake with mozzarella, then finish with basil and olive oil."], ru: ["Разогрейте духовку и противень до максимальной температуры.", "Растяните тесто и тонко распределите томаты.", "Запеките с моцареллой, затем добавьте базилик и оливковое масло."] },
    tips: { en: ["Less topping gives a better crust."], ru: ["Чем меньше начинки, тем лучше корочка."] }, tags: ["pizza", "vegetarian", "пицца", "вегетарианское"]
  },
  {
    slug: "georgian-khachapuri",
    title: { en: "Adjarian Khachapuri", ru: "Аджарский хачапури" },
    description: { en: "A warm boat of bread filled with molten cheese, egg and butter.", ru: "Тёплая лодочка из теста с расплавленным сыром, яйцом и маслом." },
    image: photo("photo-1547592180-85f173990554"), cuisine: "georgian", category: "baking", mealType: ["breakfast", "lunch"], diet: ["vegetarian"], difficulty: "hard", prepTimeMinutes: 35, cookTimeMinutes: 18, servings: 2, rating: 4.9, calories: 740,
    ingredients: { en: ["300 g pizza dough", "250 g sulguni or mozzarella", "2 eggs", "30 g butter"], ru: ["300 г теста", "250 г сулугуни или моцареллы", "2 яйца", "30 г сливочного масла"] },
    steps: { en: ["Shape dough into two boats with high edges.", "Fill with cheese and bake until deeply golden.", "Crack in eggs, return briefly to oven and finish with butter."], ru: ["Сформируйте из теста две лодочки с высокими бортиками.", "Наполните сыром и запекайте до насыщенного золотистого цвета.", "Добавьте яйца, верните в духовку на несколько минут и завершите маслом."] },
    tips: { en: ["Mix the egg and butter into the cheese at the table."], ru: ["Смешайте яйцо и масло с сыром уже за столом."] }, tags: ["bread", "cheese", "грузинская", "сыр"]
  },
  {
    slug: "chicken-ramen",
    title: { en: "Chicken Ramen", ru: "Рамен с курицей" },
    description: { en: "Comforting broth with springy noodles, chicken and jammy egg.", ru: "Согревающий бульон с упругой лапшой, курицей и яйцом с кремовым желтком." },
    image: photo("photo-1569718212165-3a8278d5f624"), cuisine: "japanese", category: "soups", mealType: ["lunch", "dinner"], diet: ["omnivore"], difficulty: "medium", prepTimeMinutes: 20, cookTimeMinutes: 35, servings: 2, rating: 4.7, calories: 590,
    ingredients: { en: ["700 ml chicken stock", "2 ramen noodle nests", "200 g chicken thigh", "2 soft-boiled eggs"], ru: ["700 мл куриного бульона", "2 порции лапши рамен", "200 г куриного бедра", "2 яйца всмятку"] },
    steps: { en: ["Simmer stock with ginger, garlic and soy sauce.", "Sear chicken and slice it finely.", "Cook noodles, divide into bowls and pour over hot broth with toppings."], ru: ["Протомите бульон с имбирём, чесноком и соевым соусом.", "Обжарьте курицу и тонко нарежьте.", "Отварите лапшу, разложите по мискам и залейте горячим бульоном с добавками."] },
    tips: { en: ["Add noodles just before serving."], ru: ["Добавляйте лапшу непосредственно перед подачей."] }, tags: ["ramen", "noodles", "рамен", "лапша"]
  },
  {
    slug: "vegetable-curry",
    title: { en: "Coconut Vegetable Curry", ru: "Овощной карри с кокосом" },
    description: { en: "A weeknight curry of vegetables, chickpeas and fragrant coconut broth.", ru: "Быстрый карри из овощей, нута и ароматного кокосового соуса." },
    image: photo("photo-1601050690597-df0568f70950"), cuisine: "indian", category: "vegetarian", mealType: ["dinner"], diet: ["vegan", "gluten-free"], difficulty: "easy", prepTimeMinutes: 15, cookTimeMinutes: 25, servings: 4, rating: 4.7, calories: 420,
    ingredients: { en: ["1 tbsp curry paste", "400 ml coconut milk", "400 g mixed vegetables", "240 g cooked chickpeas"], ru: ["1 ст. л. пасты карри", "400 мл кокосового молока", "400 г разных овощей", "240 г готового нута"] },
    steps: { en: ["Fry curry paste until fragrant.", "Add coconut milk, vegetables and chickpeas.", "Simmer until tender and serve with lime."], ru: ["Обжарьте пасту карри до аромата.", "Добавьте кокосовое молоко, овощи и нут.", "Тушите до мягкости и подавайте с лаймом."] },
    tips: { en: ["Use whatever vegetables are in season."], ru: ["Используйте любые сезонные овощи."] }, tags: ["curry", "vegan", "карри", "веганское"]
  },
  {
    slug: "caesar-salad",
    title: { en: "Classic Caesar Salad", ru: "Классический салат «Цезарь»" },
    description: { en: "Crisp romaine, garlic croutons and a sharp anchovy dressing.", ru: "Хрустящий ромэн, чесночные сухарики и пикантная заправка с анчоусами." },
    image: photo("photo-1546793665-c74683f339c1"), cuisine: "american", category: "salads", mealType: ["lunch"], diet: ["omnivore"], difficulty: "easy", prepTimeMinutes: 18, cookTimeMinutes: 7, servings: 2, rating: 4.6, calories: 460,
    ingredients: { en: ["1 romaine heart", "50 g parmesan", "2 cups bread cubes", "Caesar dressing"], ru: ["1 кочан ромэна", "50 г пармезана", "2 стакана хлебных кубиков", "Заправка «Цезарь»"] },
    steps: { en: ["Toast bread cubes with garlic and oil.", "Tear lettuce and shave parmesan.", "Toss just before serving with dressing and croutons."], ru: ["Подсушите хлебные кубики с чесноком и маслом.", "Порвите салат и нарежьте пармезан стружкой.", "Перед подачей перемешайте с заправкой и сухариками."] },
    tips: { en: ["Dress lettuce at the last possible moment."], ru: ["Заправляйте салат в последний момент."] }, tags: ["salad", "caesar", "салат", "цезарь"]
  },
  {
    slug: "ukrainian-borscht",
    title: { en: "Ukrainian Borscht", ru: "Украинский борщ" },
    description: { en: "A ruby beet soup with cabbage, beans and a spoonful of sour cream.", ru: "Рубиновый свекольный суп с капустой, фасолью и ложкой сметаны." },
    image: photo("photo-1547592166-23ac45744acd"), cuisine: "ukrainian", category: "soups", mealType: ["lunch", "dinner"], diet: ["vegetarian", "gluten-free"], difficulty: "medium", prepTimeMinutes: 25, cookTimeMinutes: 55, servings: 6, rating: 4.8, calories: 290,
    ingredients: { en: ["3 beets", "300 g cabbage", "1 can white beans", "1 litre vegetable stock"], ru: ["3 свёклы", "300 г капусты", "1 банка белой фасоли", "1 л овощного бульона"] },
    steps: { en: ["Sauté onion, carrot and beet until softened.", "Add stock, cabbage and beans and simmer.", "Balance with vinegar, dill and sour cream."], ru: ["Обжарьте лук, морковь и свёклу до мягкости.", "Добавьте бульон, капусту и фасоль, затем варите на слабом огне.", "Сбалансируйте вкус уксусом, укропом и сметаной."] },
    tips: { en: ["Borscht tastes even better the next day."], ru: ["На следующий день борщ становится ещё вкуснее."] }, tags: ["soup", "beetroot", "суп", "свёкла"]
  },
  {
    slug: "shakshuka",
    title: { en: "Shakshuka", ru: "Шакшука" },
    description: { en: "Eggs gently poached in a spiced tomato and pepper sauce.", ru: "Яйца, бережно приготовленные в пряном соусе из томатов и перца." },
    image: photo("photo-1590412200988-a436970781fa"), cuisine: "mediterranean", category: "breakfast", mealType: ["breakfast", "lunch"], diet: ["vegetarian", "gluten-free"], difficulty: "easy", prepTimeMinutes: 10, cookTimeMinutes: 20, servings: 2, rating: 4.8, calories: 340,
    ingredients: { en: ["4 eggs", "1 red pepper", "400 g tomatoes", "1 tsp cumin"], ru: ["4 яйца", "1 красный перец", "400 г томатов", "1 ч. л. зиры"] },
    steps: { en: ["Soften pepper and spices in olive oil.", "Simmer tomatoes until thick and sweet.", "Make wells, add eggs and cover until just set."], ru: ["Обжарьте перец и специи в оливковом масле.", "Тушите томаты до густоты и сладости.", "Сделайте углубления, добавьте яйца и готовьте под крышкой до схватывания белка."] },
    tips: { en: ["Serve with warm bread for the sauce."], ru: ["Подавайте с тёплым хлебом, чтобы собрать соус."] }, tags: ["eggs", "breakfast", "яйца", "завтрак"]
  },
  {
    slug: "beef-tacos",
    title: { en: "Beef Tacos", ru: "Такос с говядиной" },
    description: { en: "Bright lime beef, soft corn tortillas and a fresh tomato salsa.", ru: "Ароматная говядина с лаймом, мягкие кукурузные тортильи и свежая сальса." },
    image: photo("photo-1551504734-5ee1c4a1479b"), cuisine: "mexican", category: "quick-meals", mealType: ["dinner"], diet: ["omnivore", "gluten-free"], difficulty: "easy", prepTimeMinutes: 15, cookTimeMinutes: 15, servings: 4, rating: 4.7, calories: 510,
    ingredients: { en: ["400 g ground beef", "8 corn tortillas", "1 tsp cumin", "Tomato salsa"], ru: ["400 г говяжьего фарша", "8 кукурузных тортильй", "1 ч. л. зиры", "Томатная сальса"] },
    steps: { en: ["Brown beef with cumin, garlic and lime.", "Warm tortillas in a dry pan.", "Fill with beef, salsa and coriander."], ru: ["Обжарьте фарш с зирой, чесноком и лаймом.", "Прогрейте тортильи на сухой сковороде.", "Наполните мясом, сальсой и кинзой."] },
    tips: { en: ["Warm tortillas make the whole dish better."], ru: ["Тёплые тортильи заметно улучшают блюдо."] }, tags: ["tacos", "mexican", "такос", "мексиканская"]
  },
  {
    slug: "french-onion-soup",
    title: { en: "French Onion Soup", ru: "Французский луковый суп" },
    description: { en: "Slowly caramelised onions beneath bubbling Gruyère toast.", ru: "Долго карамелизированный лук под хрустящим тостом с расплавленным грюйером." },
    image: photo("photo-1547592180-85f173990554"), cuisine: "french", category: "soups", mealType: ["lunch", "dinner"], diet: ["vegetarian"], difficulty: "medium", prepTimeMinutes: 15, cookTimeMinutes: 55, servings: 4, rating: 4.8, calories: 410,
    ingredients: { en: ["800 g onions", "1 litre vegetable stock", "4 baguette slices", "120 g Gruyère"], ru: ["800 г лука", "1 л овощного бульона", "4 ломтика багета", "120 г грюйера"] },
    steps: { en: ["Cook onions slowly until deep amber.", "Deglaze and simmer with stock.", "Top bowls with toast and cheese, then grill."], ru: ["Томите лук на медленном огне до глубокого янтарного цвета.", "Добавьте жидкость и бульон, затем проварите.", "Накройте миски тостом и сыром, после чего запеките под грилем."] },
    tips: { en: ["Patience is the main ingredient."], ru: ["Главный ингредиент — терпение."] }, tags: ["soup", "onion", "суп", "лук"]
  },
  {
    slug: "buttermilk-pancakes",
    title: { en: "Buttermilk Pancakes", ru: "Панкейки на пахте" },
    description: { en: "Tender, tall pancakes ready for butter and berries.", ru: "Нежные высокие панкейки, созданные для масла и ягод." },
    image: photo("photo-1528207776546-365bb710ee93"), cuisine: "american", category: "breakfast", mealType: ["breakfast"], diet: ["vegetarian"], difficulty: "easy", prepTimeMinutes: 10, cookTimeMinutes: 15, servings: 4, rating: 4.7, calories: 380,
    ingredients: { en: ["200 g flour", "300 ml buttermilk", "1 egg", "1 tsp baking powder"], ru: ["200 г муки", "300 мл пахты", "1 яйцо", "1 ч. л. разрыхлителя"] },
    steps: { en: ["Whisk dry ingredients in one bowl.", "Fold in buttermilk and egg without overmixing.", "Cook small rounds until bubbles appear, then flip."], ru: ["Смешайте сухие ингредиенты в одной миске.", "Введите пахту и яйцо, не перемешивая слишком долго.", "Жарьте небольшие порции до пузырьков, затем переверните."] },
    tips: { en: ["A few lumps make pancakes lighter."], ru: ["Несколько комочков сделают панкейки пышнее."] }, tags: ["pancakes", "breakfast", "панкейки", "завтрак"]
  },
  {
    slug: "apple-pie",
    title: { en: "Rustic Apple Pie", ru: "Домашний яблочный пирог" },
    description: { en: "Buttery pastry surrounds cinnamon apples with a crisp, sugared lid.", ru: "Сливочное тесто окружает яблоки с корицей под хрустящей сахарной корочкой." },
    image: photo("photo-1562007908-17c67e878c88"), cuisine: "american", category: "desserts", mealType: ["dessert"], diet: ["vegetarian"], difficulty: "medium", prepTimeMinutes: 30, cookTimeMinutes: 45, servings: 8, rating: 4.8, calories: 440,
    ingredients: { en: ["2 pastry sheets", "800 g apples", "80 g brown sugar", "1 tsp cinnamon"], ru: ["2 пласта песочного теста", "800 г яблок", "80 г коричневого сахара", "1 ч. л. корицы"] },
    steps: { en: ["Toss sliced apples with sugar and cinnamon.", "Line a pie dish, fill and add the top crust.", "Bake until bubbling and deeply golden."], ru: ["Перемешайте ломтики яблок с сахаром и корицей.", "Выложите форму тестом, добавьте начинку и верхний слой.", "Запекайте до пузырьков и насыщенной золотистой корочки."] },
    tips: { en: ["Let the pie rest before slicing."], ru: ["Дайте пирогу отдохнуть перед нарезкой."] }, tags: ["apple", "pie", "яблоко", "пирог"]
  },
  {
    slug: "greek-salad",
    title: { en: "Greek Salad", ru: "Греческий салат" },
    description: { en: "Cucumber, tomato, olives and feta in a sharp oregano dressing.", ru: "Огурец, томаты, оливки и фета в яркой заправке с орегано." },
    image: photo("photo-1540420773420-3366772f4999"), cuisine: "mediterranean", category: "salads", mealType: ["lunch"], diet: ["vegetarian", "gluten-free"], difficulty: "easy", prepTimeMinutes: 12, cookTimeMinutes: 0, servings: 3, rating: 4.6, calories: 280,
    ingredients: { en: ["2 tomatoes", "1 cucumber", "150 g feta", "80 g Kalamata olives"], ru: ["2 томата", "1 огурец", "150 г феты", "80 г оливок каламата"] },
    steps: { en: ["Cut vegetables into generous pieces.", "Add olives and a slab of feta.", "Dress with olive oil, oregano and red wine vinegar."], ru: ["Нарежьте овощи крупными кусочками.", "Добавьте оливки и большой кусок феты.", "Заправьте оливковым маслом, орегано и красным винным уксусом."] },
    tips: { en: ["Do not crumble the feta too finely."], ru: ["Не крошите фету слишком мелко."] }, tags: ["salad", "feta", "салат", "фета"]
  },
  {
    slug: "pad-thai",
    title: { en: "Vegetable Pad Thai", ru: "Овощной пад-тай" },
    description: { en: "Chewy rice noodles with tamarind, tofu, crunchy peanuts and lime.", ru: "Упругая рисовая лапша с тамариндом, тофу, хрустящим арахисом и лаймом." },
    image: photo("photo-1559314809-0d155014e29e"), cuisine: "japanese", category: "quick-meals", mealType: ["dinner"], diet: ["vegan", "gluten-free"], difficulty: "medium", prepTimeMinutes: 20, cookTimeMinutes: 12, servings: 3, rating: 4.6, calories: 480,
    ingredients: { en: ["200 g rice noodles", "180 g tofu", "2 tbsp tamarind paste", "60 g bean sprouts"], ru: ["200 г рисовой лапши", "180 г тофу", "2 ст. л. тамариндовой пасты", "60 г ростков фасоли"] },
    steps: { en: ["Soak noodles until pliable.", "Sear tofu and build a sweet-sour tamarind sauce.", "Toss noodles quickly with sprouts, peanuts and lime."], ru: ["Замочите лапшу до мягкости.", "Обжарьте тофу и приготовьте кисло-сладкий соус с тамариндом.", "Быстро перемешайте лапшу с ростками, арахисом и лаймом."] },
    tips: { en: ["Have everything prepared before you start cooking."], ru: ["Подготовьте всё заранее — готовится блюдо очень быстро."] }, tags: ["noodles", "vegan", "лапша", "веганское"]
  },
  {
    slug: "mushroom-risotto",
    title: { en: "Mushroom Risotto", ru: "Ризотто с грибами" },
    description: { en: "Creamy arborio rice with browned mushrooms and parmesan.", ru: "Кремовый рис арборио с подрумяненными грибами и пармезаном." },
    image: photo("photo-1476124369491-e7addf5db371"), cuisine: "italian", category: "dinner", mealType: ["dinner"], diet: ["vegetarian", "gluten-free"], difficulty: "medium", prepTimeMinutes: 15, cookTimeMinutes: 35, servings: 4, rating: 4.8, calories: 510,
    ingredients: { en: ["300 g arborio rice", "300 g mushrooms", "1 litre vegetable stock", "70 g parmesan"], ru: ["300 г риса арборио", "300 г грибов", "1 л овощного бульона", "70 г пармезана"] },
    steps: { en: ["Brown mushrooms separately until concentrated.", "Toast rice, then add hot stock a ladle at a time.", "Stir in butter, parmesan and mushrooms off the heat."], ru: ["Обжарьте грибы отдельно до концентрированного вкуса.", "Прогрейте рис и добавляйте горячий бульон по половнику.", "Снимите с огня и вмешайте масло, пармезан и грибы."] },
    tips: { en: ["The final texture should flow slowly on a plate."], ru: ["Готовое ризотто должно медленно растекаться по тарелке."] }, tags: ["risotto", "mushroom", "ризотто", "грибы"]
  },
  {
    slug: "latvian-rye-bread-dessert",
    title: { en: "Latvian Rye Bread Dessert", ru: "Латвийский десерт из ржаного хлеба" },
    description: { en: "Toasted rye crumbs layered with tart berries and vanilla cream.", ru: "Обжаренные ржаные крошки слоями с терпкими ягодами и ванильным кремом." },
    image: photo("photo-1488477181946-6428a0291777"), cuisine: "latvian", category: "desserts", mealType: ["dessert"], diet: ["vegetarian"], difficulty: "easy", prepTimeMinutes: 20, cookTimeMinutes: 10, servings: 4, rating: 4.5, calories: 360,
    ingredients: { en: ["180 g rye bread crumbs", "250 ml cream", "200 g mixed berries", "2 tbsp sugar"], ru: ["180 г ржаных хлебных крошек", "250 мл сливок", "200 г ягод", "2 ст. л. сахара"] },
    steps: { en: ["Toast rye crumbs with a little sugar.", "Whip cream with vanilla.", "Layer crumbs, berries and cream in glasses."], ru: ["Подсушите ржаные крошки с небольшим количеством сахара.", "Взбейте сливки с ванилью.", "Выложите в стаканы слоями крошки, ягоды и крем."] },
    tips: { en: ["Assemble shortly before serving for the best crunch."], ru: ["Собирайте незадолго до подачи, чтобы сохранить хруст."] }, tags: ["rye", "berries", "рожь", "ягоды"]
  },
  {
    slug: "falafel-bowl",
    title: { en: "Falafel Bowl", ru: "Боул с фалафелем" },
    description: { en: "Herby chickpea falafel with creamy tahini, grains and crunchy vegetables.", ru: "Пряный фалафель из нута с кремовым тахини, крупой и хрустящими овощами." },
    image: photo("photo-1540420773420-3366772f4999"), cuisine: "mediterranean", category: "vegetarian", mealType: ["lunch", "dinner"], diet: ["vegan"], difficulty: "medium", prepTimeMinutes: 25, cookTimeMinutes: 18, servings: 4, rating: 4.7, calories: 520,
    ingredients: { en: ["300 g cooked chickpeas", "1 cup cooked bulgur", "3 tbsp tahini", "Cucumber and herbs"], ru: ["300 г готового нута", "1 стакан готового булгура", "3 ст. л. тахини", "Огурец и зелень"] },
    steps: { en: ["Blend chickpeas with herbs and spices, then shape balls.", "Bake or pan-fry until crisp.", "Build bowls with grains, vegetables and lemon tahini."], ru: ["Измельчите нут с зеленью и специями, сформируйте шарики.", "Запеките или обжарьте до хрустящей корочки.", "Соберите боулы с крупой, овощами и лимонным тахини."] },
    tips: { en: ["Dry chickpeas thoroughly for a crisp falafel."], ru: ["Хорошо обсушите нут для хрустящего фалафеля."] }, tags: ["falafel", "vegan", "фалафель", "веганское"]
  },
  {
    slug: "salmon-teriyaki",
    title: { en: "Salmon Teriyaki", ru: "Лосось терияки" },
    description: { en: "Lacquered salmon with rice, sesame and quick greens.", ru: "Глянцевый лосось с рисом, кунжутом и быстро приготовленной зеленью." },
    image: photo("photo-1467003909585-2f8a72700288"), cuisine: "japanese", category: "dinner", mealType: ["dinner"], diet: ["pescatarian", "gluten-free"], difficulty: "easy", prepTimeMinutes: 10, cookTimeMinutes: 15, servings: 2, rating: 4.8, calories: 540,
    ingredients: { en: ["2 salmon fillets", "3 tbsp tamari", "1 tbsp maple syrup", "2 bowls cooked rice"], ru: ["2 филе лосося", "3 ст. л. тамари", "1 ст. л. кленового сиропа", "2 порции готового риса"] },
    steps: { en: ["Whisk tamari, maple and ginger.", "Sear salmon skin-side down until crisp.", "Brush with sauce and cook until glossy."], ru: ["Смешайте тамари, сироп и имбирь.", "Обжарьте лосось кожей вниз до хруста.", "Смажьте соусом и готовьте до глянцевой глазури."] },
    tips: { en: ["Use a hot pan and do not move the fish too early."], ru: ["Используйте горячую сковороду и не двигайте рыбу слишком рано."] }, tags: ["salmon", "fish", "лосось", "рыба"]
  },
  {
    slug: "lentil-soup",
    title: { en: "Red Lentil Soup", ru: "Суп из красной чечевицы" },
    description: { en: "Velvety red lentils with lemon, cumin and a drizzle of olive oil.", ru: "Бархатистая красная чечевица с лимоном, зирой и каплей оливкового масла." },
    image: photo("photo-1547592180-85f173990554"), cuisine: "mediterranean", category: "soups", mealType: ["lunch"], diet: ["vegan", "gluten-free"], difficulty: "easy", prepTimeMinutes: 10, cookTimeMinutes: 30, servings: 5, rating: 4.7, calories: 320,
    ingredients: { en: ["250 g red lentils", "1 carrot", "1 tsp cumin", "1 litre vegetable stock"], ru: ["250 г красной чечевицы", "1 морковь", "1 ч. л. зиры", "1 л овощного бульона"] },
    steps: { en: ["Sauté carrot, onion and cumin.", "Add lentils and stock, then simmer until soft.", "Blend partly and brighten with lemon."], ru: ["Обжарьте морковь, лук и зиру.", "Добавьте чечевицу и бульон, варите до мягкости.", "Частично пробейте блендером и добавьте лимон."] },
    tips: { en: ["A spoon of chilli oil adds welcome heat."], ru: ["Ложка масла с чили добавит приятную остроту."] }, tags: ["lentils", "vegan", "чечевица", "веганское"]
  },
  {
    slug: "tiramisu",
    title: { en: "Classic Tiramisu", ru: "Классический тирамису" },
    description: { en: "Coffee-soaked savoiardi under a cloud of mascarpone cream.", ru: "Савоярди, пропитанные кофе, под облаком крема из маскарпоне." },
    image: photo("photo-1571877227200-a0d98ea607e9"), cuisine: "italian", category: "desserts", mealType: ["dessert"], diet: ["vegetarian"], difficulty: "medium", prepTimeMinutes: 25, cookTimeMinutes: 0, servings: 8, rating: 4.9, calories: 450,
    ingredients: { en: ["250 g mascarpone", "24 savoiardi", "250 ml espresso", "Cocoa powder"], ru: ["250 г маскарпоне", "24 печенья савоярди", "250 мл эспрессо", "Какао-порошок"] },
    steps: { en: ["Whisk mascarpone into a light cream.", "Dip savoiardi briefly in cooled coffee.", "Layer with cream and chill before dusting with cocoa."], ru: ["Взбейте маскарпоне в лёгкий крем.", "Быстро окунайте савоярди в остывший кофе.", "Выложите слоями с кремом, охладите и посыпьте какао."] },
    tips: { en: ["An overnight rest gives the cleanest slices."], ru: ["После ночи в холодильнике порции будут аккуратнее."] }, tags: ["tiramisu", "coffee", "тирамису", "кофе"]
  },
  {
    slug: "ratatouille",
    title: { en: "Summer Ratatouille", ru: "Летний рататуй" },
    description: { en: "A gentle Provençal stew of tomato, eggplant, courgette and herbs.", ru: "Нежное провансальское рагу из томатов, баклажанов, кабачков и трав." },
    image: photo("photo-1476224203421-9ac39bcb3327"), cuisine: "french", category: "vegetarian", mealType: ["lunch", "dinner"], diet: ["vegan", "gluten-free"], difficulty: "easy", prepTimeMinutes: 20, cookTimeMinutes: 40, servings: 4, rating: 4.6, calories: 230,
    ingredients: { en: ["1 eggplant", "2 courgettes", "4 tomatoes", "Herbes de Provence"], ru: ["1 баклажан", "2 кабачка", "4 томата", "Прованские травы"] },
    steps: { en: ["Salt eggplant briefly, then pat dry.", "Cook vegetables in stages so they keep their shape.", "Simmer together with tomato and herbs until glossy."], ru: ["Посолите баклажан ненадолго, затем обсушите.", "Готовьте овощи поэтапно, чтобы они сохранили форму.", "Тушите вместе с томатами и травами до глянцевого соуса."] },
    tips: { en: ["Serve warm or at room temperature."], ru: ["Подавайте тёплым или комнатной температуры."] }, tags: ["vegetables", "vegan", "овощи", "веганское"]
  },
  {
    slug: "hummus-plate",
    title: { en: "Creamy Hummus Plate", ru: "Тарелка с кремовым хумусом" },
    description: { en: "Ultra-smooth chickpea hummus with warm flatbread and crisp vegetables.", ru: "Нежный хумус из нута с тёплой лепёшкой и хрустящими овощами." },
    image: photo("photo-1577805947697-89e18249d767"), cuisine: "mediterranean", category: "quick-meals", mealType: ["lunch"], diet: ["vegan"], difficulty: "easy", prepTimeMinutes: 12, cookTimeMinutes: 0, servings: 4, rating: 4.6, calories: 330,
    ingredients: { en: ["400 g cooked chickpeas", "3 tbsp tahini", "1 lemon", "1 garlic clove"], ru: ["400 г готового нута", "3 ст. л. тахини", "1 лимон", "1 зубчик чеснока"] },
    steps: { en: ["Blend chickpeas until very smooth.", "Add tahini, lemon, garlic and ice water.", "Spread on a plate and finish with olive oil."], ru: ["Измельчите нут до очень гладкой текстуры.", "Добавьте тахини, лимон, чеснок и ледяную воду.", "Выложите на тарелку и завершите оливковым маслом."] },
    tips: { en: ["Peeling chickpeas gives an extra-silky result."], ru: ["Очищенный нут сделает хумус особенно шелковистым."] }, tags: ["hummus", "vegan", "хумус", "веганское"]
  },
  {
    slug: "new-york-cheesecake",
    title: { en: "New York Cheesecake", ru: "Нью-йоркский чизкейк" },
    description: { en: "A tall, tangy baked cheesecake with a buttery biscuit base.", ru: "Высокий запечённый чизкейк с лёгкой кислинкой и масляной основой из печенья." },
    image: photo("photo-1533134242443-d4fd215305ad"), cuisine: "american", category: "desserts", mealType: ["dessert"], diet: ["vegetarian"], difficulty: "hard", prepTimeMinutes: 30, cookTimeMinutes: 60, servings: 10, rating: 4.8, calories: 530,
    ingredients: { en: ["250 g digestive biscuits", "600 g cream cheese", "180 g sugar", "3 eggs"], ru: ["250 г песочного печенья", "600 г сливочного сыра", "180 г сахара", "3 яйца"] },
    steps: { en: ["Press biscuit crumbs into a lined tin.", "Mix cream cheese, sugar and eggs until just smooth.", "Bake gently, cool slowly and chill overnight."], ru: ["Утрамбуйте крошку печенья в подготовленной форме.", "Смешайте сливочный сыр, сахар и яйца до однородности.", "Запекайте бережно, медленно охладите и оставьте в холодильнике на ночь."] },
    tips: { en: ["Avoid overbeating to prevent cracks."], ru: ["Не взбивайте слишком долго, чтобы не появились трещины."] }, tags: ["cheesecake", "dessert", "чизкейк", "десерт"]
  },
  {
    slug: "berry-smoothie",
    title: { en: "Berry Oat Smoothie", ru: "Ягодный смузи с овсянкой" },
    description: { en: "A frosty berry smoothie with oats, yogurt and a hint of vanilla.", ru: "Холодный ягодный смузи с овсянкой, йогуртом и ноткой ванили." },
    image: photo("photo-1553530666-ba11a90a0868"), cuisine: "american", category: "drinks", mealType: ["breakfast"], diet: ["vegetarian", "gluten-free"], difficulty: "easy", prepTimeMinutes: 5, cookTimeMinutes: 0, servings: 2, rating: 4.5, calories: 290,
    ingredients: { en: ["250 g frozen berries", "200 ml yogurt", "40 g oats", "1 banana"], ru: ["250 г замороженных ягод", "200 мл йогурта", "40 г овсянки", "1 банан"] },
    steps: { en: ["Add all ingredients to a high-speed blender.", "Blend until creamy and bright.", "Serve immediately with extra berries."], ru: ["Сложите все ингредиенты в мощный блендер.", "Измельчите до кремовой яркой текстуры.", "Сразу подавайте, добавив сверху ягоды."] },
    tips: { en: ["Freeze the banana for an extra-thick texture."], ru: ["Заморозьте банан для особенно густой текстуры."] }, tags: ["smoothie", "berries", "смузи", "ягоды"]
  }
];

export const recipes = seeds.map(makeRecipe);

export const recipeBySlug = (slug: string) => recipes.find((recipe) => recipe.slug === slug);
