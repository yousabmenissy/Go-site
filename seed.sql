INSERT INTO books (
        "title",
        "author",
        "release_date",
        "publisher",
        "country",
        "isbn",
        "pages",
        "language",
        "description",
        "genres",
        "cover_image_url",
        "price"
    )
VALUES (
        'animal farm',
        'george orwell',
        '1945-8-17',
        'Secker and Warburg',
        'UK',
        '978-9386690364',
        92,
        'en',
        'This is a classic tale of humanity awash in totalitarianism. 
A farm is taken over by its overworked, mistreated animals. 
With flaming idealism and stirring slogans, they set out to create a paradise of progress, justice, and equality. 
First published during the epoch of Stalinist Russia, today it is clear that wherever and whenever freedom is attacked, and under whatever banner, the cutting clarity and savage comedy of Orwell''s masterpiece is a message still ferociously fresh.',
        ARRAY ['political satire',
        'fiction'],
        'https://upload.wikimedia.org/wikipedia/commons/f/fb/Animal_Farm_-_1st_edition.jpg',
        9.99
    );
-- INSERT INTO books ("title", "author")
-- VALUES ('the hobbit', 'J.R.R Tolkein');
-- INSERT INTO books ("title", "author")
-- VALUES ('the lord of the rings', 'J.R.R Tolkein');
/*
 id BIGSERIAL PRIMARY KEY,
 title TEXT NOT NULL,
 subtitle TEXT NOT NULL,
 author TEXT NOT NULL,
 release_date DATE NOT NULL,
 publisher TEXT NOT NULL,
 isbn TEXT NOT NULL,
 pages INT NOT NULL,
 language TEXT NOT NULL,
 description TEXT NOT NULL,
 categories TEXT [] NOT NULL,
 cover_image_url TEXT NOT NULL,
 price NUMERIC(10, 2)
 */