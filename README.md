# js packagemanager

This is a go module to deal with javascript ecosystem packagemanager.

# features
- detection of common package managers such as yarn, npm, pnpm.
- get workspaces package.json when dealing with mono[repo|space]
- load a package.json into a struct (provided by the packageJson module in case you only need this)

Other features will come soon like some common commands to launch on the sytem with thoose package managers
A better documentation will also come soon


## history of this package
The code in this repository was originally inspired by https://github.com/replit/upm and largely modified by the turbo team at vercel
This derivative work is not affiliated to any of thoose, but it is important to know the origin of the code in this repository.

This package was then extracted from turbo and cleaned up to be usable in other contexts than turbo repo.
This we lost some capabilities proposed by the original code, like pruning lock files, that will perhaps be re-integrated in the future, but there's no particular plan on this. I hope this will be helpful for a bunch of crazy devs around and as always you can propose PR to make this a better tool.

