# rex

Experimental relational NoSQL database. It is my playground for ideas and API will change over time. There is a lot more to do before it can be even considered interesting.

Rex has two components. Package table that provides table structure and table operations and command line interface.

## Command rex

``` shell
% rex union -ts movies.json
┏━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━┓
┃ director         │ title                          │ year ┃
┠──────────────────┼────────────────────────────────┼──────┨
┃ Denis Villeneuve │ Blade Runner 2049              │ 2017 ┃
┃ Denis Villeneuve │ Dune                           │ 2021 ┃
┃ James Gunn       │ Guardians of the Galaxy Vol. 3 │ 2023 ┃
┃ Denis Villeneuve │ Dune: Part Two                 │ 2024 ┃
┗━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━┛

% rex natural-join -ts movies.json -ia '[{"year": 2024}]'
┏━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━┯━━━━━━┓
┃ director         │ title          │ year ┃
┠──────────────────┼────────────────┼──────┨
┃ Denis Villeneuve │ Dune: Part Two │ 2024 ┃
┗━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━┷━━━━━━┛
```
