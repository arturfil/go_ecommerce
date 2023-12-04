INSERT INTO `sessions` 
    (name, description, price, session_date, image, created_at, updated_at)
    VALUES ("Calculus 106", "Integrals and Derivatives mathematics", 8000, "2023-11-29", "image.jpg", NOW(), NOW());

INSERT INTO `transaction_status`
    (name)
    VALUES
        ("Pending"),
        ("Cleared"),
        ("Declined"),
        ("Refunded"),
        ("Partially refunded")
    ;
