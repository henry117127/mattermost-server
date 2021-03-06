// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package config

import "github.com/pkg/errors"

func Migrate(from, to string) error {
	source, err := NewStore(from, false)
	if err != nil {
		return errors.Wrapf(err, "failed to access source config %s", from)
	}

	destination, err := NewStore(to, false)
	if err != nil {
		return errors.Wrapf(err, "failed to access destination config %s", to)
	}

	sourceConfig := source.Get()
	if _, err = destination.Set(sourceConfig); err != nil {
		return errors.Wrapf(err, "failed to set config")
	}

	files := []string{*sourceConfig.SamlSettings.IdpCertificateFile, *sourceConfig.SamlSettings.PublicCertificateFile,
		*sourceConfig.SamlSettings.PrivateKeyFile}

	for _, file := range files {
		err = migrateFile(file, source, destination)

		if err != nil {
			return err
		}
	}
	return nil
}

func migrateFile(name string, source Store, destination Store) error {
	fileExists, err := source.HasFile(name)
	if err != nil {
		return errors.Wrapf(err, "failed to check existence of %s", name)
	}

	if fileExists {
		file, err := source.GetFile(name)
		err = destination.SetFile(name, file)
		if err != nil {
			return errors.Wrapf(err, "failed to migrate %s", name)
		}
	}

	return nil
}
