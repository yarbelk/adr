package repo

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"gitlab.com/yarbelk/adr/src/adr"
	"gitlab.com/yarbelk/adr/src/ioutils"
	"gitlab.com/yarbelk/adr/src/serializer"
)

type FileRepo struct {
	dir  string
	adrs map[int]adr.ADR
	keys []int

	newUnmarshaller func(io.Reader) serializer.Unmarshaller
	newMarshaller   func(io.Writer) serializer.Marshaller
}

// NewFileRepo with provided unmarshalers and marshallers.  Instantiates from the directory
func NewFileRepo(dir string, newUnmarshaller func(io.Reader) serializer.Unmarshaller, newMarshaller func(io.Writer) serializer.Marshaller) FileRepo {
	files, err := ioutils.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	keys := make([]int, 0, len(files))
	mapADR := make(map[int]adr.ADR)
	re := regexp.MustCompile("^ADR-[[:digit:]]{4}$")
	for _, file := range files {
		if !re.MatchString(file.Name()) {
			continue
		}
		fp := filepath.Join(dir, file.Name())
		newADR := new(adr.ADR)

		// Sooo.  this is the make the defer fire earlier; because i
		// don't want them all to pile up.  This is *not* supposed to
		// be a go routine.
		func() {
			f, err := os.OpenFile(fp, os.O_RDONLY, 0644)
			if err != nil {
				log.Fatal("can't read the ADRs in the dir", err)
			}
			defer f.Close()

			if err := newUnmarshaller(f).Unmarshal(newADR); err != nil {
				log.Fatal("failed to load old ADR in update related", err)
			}
			mapADR[newADR.Number] = *newADR
			keys = append(keys, newADR.Number)
			return
		}()
	}
	sort.Ints(keys)
	repo := FileRepo{
		dir:             dir,
		keys:            keys,
		adrs:            mapADR,
		newMarshaller:   newMarshaller,
		newUnmarshaller: newUnmarshaller,
	}

	return repo
}

// List out the adrs in order
func (repo FileRepo) List() []*adr.ADR {
	results := make([]*adr.ADR, 0, len(repo.keys))
	for _, key := range repo.keys {
		a, err := repo.Get(key)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, a)
	}
	return results
}

// Get an adr by its id number (so, ADR-0001, pass in a '1')
// will return an error if it can't be found or on deserializetion
// problems
func (repo FileRepo) Get(n int) (*adr.ADR, error) {
	adr, ok := repo.adrs[n]
	if !ok {
		return nil, fmt.Errorf("Cannot find ADR-%04d", n)
	}

	return &adr, nil
}

// Add a new adr to the repo.  If one with that number already exists
// then  it will return an error
func (repo FileRepo) Add(newADR adr.ADR) error {
	if _, ok := repo.adrs[newADR.Number]; ok {
		return fmt.Errorf("Already have an adr %s", newADR.Filename())
	}
	repo.adrs[newADR.Number] = newADR
	fp := filepath.Join(repo.dir, newADR.Filename())
	f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return repo.newMarshaller(f).Marshal(&newADR)
}

// Update one or more pre-exisitng adrs.  They will be
// rewritten with their new data in whatever storage mechanism
// is defined.  This will use the Number to find/update them
// this will return an error if one of the adrs does not exist
// It will not update anything if one does not exist
func (repo FileRepo) Update(adrs ...adr.ADR) error {
	for _, a := range adrs {
		if _, ok := repo.adrs[a.Number]; !ok {
			return fmt.Errorf("no adr to update with number %s", a.Filename())
		}
	}

	for _, a := range adrs {
		repo.adrs[a.Number] = a
		fp := filepath.Join(repo.dir, a.Filename())
		f, err := os.OpenFile(fp, os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		return repo.newMarshaller(f).Marshal(&a)
	}
	return nil
}
