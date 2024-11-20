--
-- PostgreSQL database dump
--

-- Dumped from database version 17.1 (Debian 17.1-1.pgdg120+1)
-- Dumped by pg_dump version 17.1 (Debian 17.1-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: dish_images; Type: TABLE; Schema: public; Owner: dancingponysvc
--

CREATE TABLE public.dish_images (
    id integer NOT NULL,
    dish_id integer NOT NULL,
    filename character varying(255) NOT NULL,
    content bytea NOT NULL,
    uploaded_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.dish_images OWNER TO dancingponysvc;

--
-- Name: dish_images_id_seq; Type: SEQUENCE; Schema: public; Owner: dancingponysvc
--

CREATE SEQUENCE public.dish_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dish_images_id_seq OWNER TO dancingponysvc;

--
-- Name: dish_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dancingponysvc
--

ALTER SEQUENCE public.dish_images_id_seq OWNED BY public.dish_images.id;


--
-- Name: restaurant_dishes; Type: TABLE; Schema: public; Owner: dancingponysvc
--

CREATE TABLE public.restaurant_dishes (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    price numeric(10,2) NOT NULL,
    category character varying(100),
    is_vegetarian boolean DEFAULT false,
    is_available boolean DEFAULT true,
    rating numeric(3,2),
    restaurant_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT restaurant_dishes_rating_check CHECK (((rating >= 1.00) AND (rating <= 5.00)))
);


ALTER TABLE public.restaurant_dishes OWNER TO dancingponysvc;

--
-- Name: restaurant_dishes_id_seq; Type: SEQUENCE; Schema: public; Owner: dancingponysvc
--

CREATE SEQUENCE public.restaurant_dishes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.restaurant_dishes_id_seq OWNER TO dancingponysvc;

--
-- Name: restaurant_dishes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dancingponysvc
--

ALTER SEQUENCE public.restaurant_dishes_id_seq OWNED BY public.restaurant_dishes.id;


--
-- Name: restaurant_reviews; Type: TABLE; Schema: public; Owner: dancingponysvc
--

CREATE TABLE public.restaurant_reviews (
    id integer NOT NULL,
    restaurant_id integer NOT NULL,
    user_id integer NOT NULL,
    review text NOT NULL,
    rating numeric(3,2),
    sentiment_score numeric,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT restaurant_reviews_rating_check CHECK (((rating >= 1.00) AND (rating <= 5.00)))
);


ALTER TABLE public.restaurant_reviews OWNER TO dancingponysvc;

--
-- Name: restaurant_reviews_id_seq; Type: SEQUENCE; Schema: public; Owner: dancingponysvc
--

CREATE SEQUENCE public.restaurant_reviews_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.restaurant_reviews_id_seq OWNER TO dancingponysvc;

--
-- Name: restaurant_reviews_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dancingponysvc
--

ALTER SEQUENCE public.restaurant_reviews_id_seq OWNED BY public.restaurant_reviews.id;


--
-- Name: restaurants; Type: TABLE; Schema: public; Owner: dancingponysvc
--

CREATE TABLE public.restaurants (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    path_name character varying(255) NOT NULL,
    location character varying(255)
);


ALTER TABLE public.restaurants OWNER TO dancingponysvc;

--
-- Name: restaurants_id_seq; Type: SEQUENCE; Schema: public; Owner: dancingponysvc
--

CREATE SEQUENCE public.restaurants_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.restaurants_id_seq OWNER TO dancingponysvc;

--
-- Name: restaurants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dancingponysvc
--

ALTER SEQUENCE public.restaurants_id_seq OWNED BY public.restaurants.id;


--
-- Name: user_dish_ratings; Type: TABLE; Schema: public; Owner: dancingponysvc
--

CREATE TABLE public.user_dish_ratings (
    id integer NOT NULL,
    restaurant_id integer NOT NULL,
    dish_id integer NOT NULL,
    user_id integer NOT NULL,
    rating numeric(3,2),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_dish_ratings_rating_check CHECK (((rating >= 1.00) AND (rating <= 5.00)))
);


ALTER TABLE public.user_dish_ratings OWNER TO dancingponysvc;

--
-- Name: user_dish_ratings_id_seq; Type: SEQUENCE; Schema: public; Owner: dancingponysvc
--

CREATE SEQUENCE public.user_dish_ratings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_dish_ratings_id_seq OWNER TO dancingponysvc;

--
-- Name: user_dish_ratings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dancingponysvc
--

ALTER SEQUENCE public.user_dish_ratings_id_seq OWNED BY public.user_dish_ratings.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: dancingponysvc
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    role character varying(50) NOT NULL,
    restaurant_id integer,
    hashed_password text NOT NULL,
    salt text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO dancingponysvc;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: dancingponysvc
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO dancingponysvc;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dancingponysvc
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: dish_images id; Type: DEFAULT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.dish_images ALTER COLUMN id SET DEFAULT nextval('public.dish_images_id_seq'::regclass);


--
-- Name: restaurant_dishes id; Type: DEFAULT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurant_dishes ALTER COLUMN id SET DEFAULT nextval('public.restaurant_dishes_id_seq'::regclass);


--
-- Name: restaurant_reviews id; Type: DEFAULT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurant_reviews ALTER COLUMN id SET DEFAULT nextval('public.restaurant_reviews_id_seq'::regclass);


--
-- Name: restaurants id; Type: DEFAULT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurants ALTER COLUMN id SET DEFAULT nextval('public.restaurants_id_seq'::regclass);


--
-- Name: user_dish_ratings id; Type: DEFAULT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.user_dish_ratings ALTER COLUMN id SET DEFAULT nextval('public.user_dish_ratings_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: dish_images dish_images_pkey; Type: CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.dish_images
    ADD CONSTRAINT dish_images_pkey PRIMARY KEY (id);


--
-- Name: restaurant_dishes restaurant_dishes_pkey; Type: CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurant_dishes
    ADD CONSTRAINT restaurant_dishes_pkey PRIMARY KEY (id);


--
-- Name: restaurant_reviews restaurant_reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurant_reviews
    ADD CONSTRAINT restaurant_reviews_pkey PRIMARY KEY (id);


--
-- Name: restaurants restaurants_pkey; Type: CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurants
    ADD CONSTRAINT restaurants_pkey PRIMARY KEY (id);


--
-- Name: user_dish_ratings user_dish_ratings_pkey; Type: CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.user_dish_ratings
    ADD CONSTRAINT user_dish_ratings_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: dish_images dish_images_dish_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.dish_images
    ADD CONSTRAINT dish_images_dish_id_fkey FOREIGN KEY (dish_id) REFERENCES public.restaurant_dishes(id) ON DELETE CASCADE;


--
-- Name: restaurant_reviews restaurant_reviews_restaurant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurant_reviews
    ADD CONSTRAINT restaurant_reviews_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE;


--
-- Name: restaurant_reviews restaurant_reviews_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.restaurant_reviews
    ADD CONSTRAINT restaurant_reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: user_dish_ratings user_dish_ratings_dish_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.user_dish_ratings
    ADD CONSTRAINT user_dish_ratings_dish_id_fkey FOREIGN KEY (dish_id) REFERENCES public.restaurant_dishes(id) ON DELETE CASCADE;


--
-- Name: user_dish_ratings user_dish_ratings_restaurant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.user_dish_ratings
    ADD CONSTRAINT user_dish_ratings_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE;


--
-- Name: user_dish_ratings user_dish_ratings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dancingponysvc
--

ALTER TABLE ONLY public.user_dish_ratings
    ADD CONSTRAINT user_dish_ratings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

