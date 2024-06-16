BRANCH=$1

# Clone the repository and apply patches
if [ -d "code-$BRANCH" ]; then
  cd ./code-$BRANCH || exit
  git stash
  git checkout $BRANCH
  cd ..
else
  git clone --depth 1 --branch $BRANCH https://github.com/supabase/supabase.git ./code-$BRANCH || exit
fi

for file in ./patches/*.patch; do
  patch -d ./code-$BRANCH -p1 < $file
done

find ./files -type f | while read -r file; do
  dest="./code-$BRANCH/${file#./files/}"
  mkdir -p "$(dirname "$dest")"
  cp "$file" "$dest"
done