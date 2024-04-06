# Spidey v2

## Usage

1. Clone this repository

```sh
git clone https://github.com/anirudhsudhir/Spidey-v2.git
cd spidey
```

2. Create a "seeds.txt" and add the seed links in quotes consecutively

    Sample seeds.txt

```text
"http://example.com"
"https://abcd.com"
```

3. Build the project and run Spidey.

   Pass the total allowed runtime time of the crawler in seconds
   and maximum request time per link in milliseconds as arguments.

   Do not pass decimal values as arguments

```sh
go build
./spidey 5 1 
# Here, 5s is the total runtime and 1ms is the maximum request time per link
```
