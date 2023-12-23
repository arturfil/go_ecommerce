INSERT INTO `meetings` 
    (name, description, plan_id, is_recurring, price, meeting_date, image, created_at, updated_at)
    VALUES 
        ("Calculus 106", "Integrals and Derivatives mathematics", "", 0, 8000, "2023-11-29", "image.jpg", NOW(), NOW()),
        ("Bronze Plan", "Three Weekly Sessions for the price of 2", "price_1OKUsGIAKhlfjLC86LIgQUMK", 1, 18000, "2023-11-29", "image.jpg", NOW(), NOW())
;

INSERT INTO `transaction_status`
    (name)
    VALUES
        ("Pending"),
        ("Cleared"),
        ("Declined"),
        ("Refunded"),
        ("Partially refunded")
    ;

INSERT INTO `users` 
    (first_name, last_name, email, password) 
    VALUES ('Admin','User','admin@example.com', '$2a$12$VR1wDmweaF3ZTVgEHiJrNOSi8VcS4j0eamr96A/7iOe8vlum3O3/q');
