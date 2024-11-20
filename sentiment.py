from transformers import pipeline
import psycopg2

# Sentiment analysis pipeline
sentiment_model = pipeline("sentiment-analysis")

# Database connection
conn = psycopg2.connect(
    dbname="dancingpony",
    user="dancingponysvc",
    password="password",
    host="localhost",
    port="5432"
)
cursor = conn.cursor()

# Fetch reviews
cursor.execute("SELECT id, review FROM restaurant_reviews")
reviews = cursor.fetchall()

# Analyze sentiment and update scores
for review_id, review_text in reviews:
    result = sentiment_model(review_text)
    sentiment_score = result[0]['score'] if result[0]['label'] == 'POSITIVE' else -result[0]['score']

    # Update database with sentiment score
    cursor.execute(
        "UPDATE restaurant_reviews SET sentiment_score = %s, updated_at = CURRENT_TIMESTAMP WHERE id = %s",
        (sentiment_score, review_id)
    )

conn.commit()
cursor.close()
conn.close()
