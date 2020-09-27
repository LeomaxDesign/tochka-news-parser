--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3
-- Dumped by pg_dump version 12.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: news_feeds; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.news_feeds (
    id bigint NOT NULL,
    url character varying NOT NULL,
    title character varying NOT NULL,
    type smallint NOT NULL,
    frequency smallint DEFAULT 60 NOT NULL,
    parse_count smallint DEFAULT 3 NOT NULL,
    item_tag character varying,
    title_tag character varying,
    description_tag character varying,
    link_tag character varying,
    published_tag character varying,
    img_tag character varying
);


ALTER TABLE public.news_feeds OWNER TO postgres;

--
-- Name: COLUMN news_feeds.frequency; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.news_feeds.frequency IS 'in seconds';


--
-- Name: COLUMN news_feeds.parse_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.news_feeds.parse_count IS 'news parse count';


--
-- Name: feed_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.feed_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.feed_id_seq OWNER TO postgres;

--
-- Name: feed_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.feed_id_seq OWNED BY public.news_feeds.id;


--
-- Name: news; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.news (
    id bigint NOT NULL,
    feed_id bigint,
    title character varying,
    description character varying,
    link character varying,
    published timestamp with time zone,
    parsed timestamp with time zone,
    img character varying
);


ALTER TABLE public.news OWNER TO postgres;

--
-- Name: feed_news_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.feed_news_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.feed_news_id_seq OWNER TO postgres;

--
-- Name: feed_news_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.feed_news_id_seq OWNED BY public.news.id;


--
-- Name: news id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news ALTER COLUMN id SET DEFAULT nextval('public.feed_news_id_seq'::regclass);


--
-- Name: news_feeds id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news_feeds ALTER COLUMN id SET DEFAULT nextval('public.feed_id_seq'::regclass);


--
-- Name: news feed_news_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news
    ADD CONSTRAINT feed_news_pkey PRIMARY KEY (id);


--
-- Name: news feed_news_title_description_link_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news
    ADD CONSTRAINT feed_news_title_description_link_key UNIQUE (title) INCLUDE (description, link);


--
-- Name: news_feeds feed_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news_feeds
    ADD CONSTRAINT feed_pkey PRIMARY KEY (id);


--
-- Name: news_feeds feed_url_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news_feeds
    ADD CONSTRAINT feed_url_key UNIQUE (url);


--
-- PostgreSQL database dump complete
--

