User Specification
==================

Reservoir should be user-centric; it's much easier that way.
However, there should still be a notion of teams, given that there may be multiple people committing to a set of projects.

As such, I believe that users and teams would be subtypes of owners.
Owners are those that own repositories.

This has an effect such that it limits the collaboration aspect of a repository, but we could remedy that later.

    type Owner interface {
      GetName() string
      GetEmail() string
      GetRepositories() []string
    }

GetRepositories() is a string for easier parsing.
Clients can always call the Repositories API to get information about repositories.
