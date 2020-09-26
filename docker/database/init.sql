--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3
-- Dumped by pg_dump version 12.3

-- Started on 2020-09-26 19:29:31

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
-- TOC entry 203 (class 1259 OID 33030)
-- Name: feed; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.feed (
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


ALTER TABLE public.feed OWNER TO postgres;

--
-- TOC entry 2838 (class 0 OID 0)
-- Dependencies: 203
-- Name: COLUMN feed.frequency; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.feed.frequency IS 'in seconds';


--
-- TOC entry 2839 (class 0 OID 0)
-- Dependencies: 203
-- Name: COLUMN feed.parse_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.feed.parse_count IS 'news parse count';


--
-- TOC entry 202 (class 1259 OID 33028)
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
-- TOC entry 2840 (class 0 OID 0)
-- Dependencies: 202
-- Name: feed_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.feed_id_seq OWNED BY public.feed.id;


--
-- TOC entry 205 (class 1259 OID 33042)
-- Name: feed_news; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.feed_news (
    id bigint NOT NULL,
    feed_id bigint,
    title character varying,
    description character varying,
    link character varying,
    published timestamp with time zone,
    parsed timestamp with time zone,
    img character varying
);


ALTER TABLE public.feed_news OWNER TO postgres;

--
-- TOC entry 204 (class 1259 OID 33040)
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
-- TOC entry 2841 (class 0 OID 0)
-- Dependencies: 204
-- Name: feed_news_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.feed_news_id_seq OWNED BY public.feed_news.id;


--
-- TOC entry 2695 (class 2604 OID 33033)
-- Name: feed id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feed ALTER COLUMN id SET DEFAULT nextval('public.feed_id_seq'::regclass);


--
-- TOC entry 2698 (class 2604 OID 33045)
-- Name: feed_news id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feed_news ALTER COLUMN id SET DEFAULT nextval('public.feed_news_id_seq'::regclass);


--
-- TOC entry 2704 (class 2606 OID 33050)
-- Name: feed_news feed_news_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feed_news
    ADD CONSTRAINT feed_news_pkey PRIMARY KEY (id);


--
-- TOC entry 2706 (class 2606 OID 33064)
-- Name: feed_news feed_news_title_description_link_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feed_news
    ADD CONSTRAINT feed_news_title_description_link_key UNIQUE (title) INCLUDE (description, link);


--
-- TOC entry 2700 (class 2606 OID 33039)
-- Name: feed feed_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feed
    ADD CONSTRAINT feed_pkey PRIMARY KEY (id);


--
-- TOC entry 2702 (class 2606 OID 33060)
-- Name: feed feed_url_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feed
    ADD CONSTRAINT feed_url_key UNIQUE (url);


-- Completed on 2020-09-26 19:29:31

--
-- PostgreSQL database dump complete
--

