package test

// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

import (
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/Khan/genqlient/internal/testutil"
)

// SimpleInlineFragmentRandomItemArticle includes the requested fields of the GraphQL type Article.
type SimpleInlineFragmentRandomItemArticle struct {
	Typename string `json:"__typename"`
	// ID is the identifier of the content.
	Id   testutil.ID `json:"id"`
	Name string      `json:"name"`
	Text string      `json:"text"`
}

// SimpleInlineFragmentRandomItemContent includes the requested fields of the GraphQL interface Content.
//
// SimpleInlineFragmentRandomItemContent is implemented by the following types:
// SimpleInlineFragmentRandomItemArticle
// SimpleInlineFragmentRandomItemVideo
// SimpleInlineFragmentRandomItemTopic
//
// The GraphQL type's documentation follows.
//
// Content is implemented by various types like Article, Video, and Topic.
type SimpleInlineFragmentRandomItemContent interface {
	implementsGraphQLInterfaceSimpleInlineFragmentRandomItemContent()
	// GetTypename returns the receiver's concrete GraphQL type-name (see interface doc for possible values).
	GetTypename() string
	// GetId returns the interface-field "id" from its implementation.
	// The GraphQL interface field's documentation follows.
	//
	// ID is the identifier of the content.
	GetId() testutil.ID
	// GetName returns the interface-field "name" from its implementation.
	GetName() string
}

func (v *SimpleInlineFragmentRandomItemArticle) implementsGraphQLInterfaceSimpleInlineFragmentRandomItemContent() {
}

// GetTypename is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemArticle) GetTypename() string { return v.Typename }

// GetId is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemArticle) GetId() testutil.ID { return v.Id }

// GetName is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemArticle) GetName() string { return v.Name }

func (v *SimpleInlineFragmentRandomItemVideo) implementsGraphQLInterfaceSimpleInlineFragmentRandomItemContent() {
}

// GetTypename is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemVideo) GetTypename() string { return v.Typename }

// GetId is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemVideo) GetId() testutil.ID { return v.Id }

// GetName is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemVideo) GetName() string { return v.Name }

func (v *SimpleInlineFragmentRandomItemTopic) implementsGraphQLInterfaceSimpleInlineFragmentRandomItemContent() {
}

// GetTypename is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemTopic) GetTypename() string { return v.Typename }

// GetId is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemTopic) GetId() testutil.ID { return v.Id }

// GetName is a part of, and documented with, the interface SimpleInlineFragmentRandomItemContent.
func (v *SimpleInlineFragmentRandomItemTopic) GetName() string { return v.Name }

func __unmarshalSimpleInlineFragmentRandomItemContent(v *SimpleInlineFragmentRandomItemContent, m json.RawMessage) error {
	if string(m) == "null" {
		return nil
	}

	var tn struct {
		TypeName string `json:"__typename"`
	}
	err := json.Unmarshal(m, &tn)
	if err != nil {
		return err
	}

	switch tn.TypeName {
	case "Article":
		*v = new(SimpleInlineFragmentRandomItemArticle)
		return json.Unmarshal(m, *v)
	case "Video":
		*v = new(SimpleInlineFragmentRandomItemVideo)
		return json.Unmarshal(m, *v)
	case "Topic":
		*v = new(SimpleInlineFragmentRandomItemTopic)
		return json.Unmarshal(m, *v)
	case "":
		return fmt.Errorf(
			"Response was missing Content.__typename")
	default:
		return fmt.Errorf(
			`Unexpected concrete type for SimpleInlineFragmentRandomItemContent: "%v"`, tn.TypeName)
	}
}

// SimpleInlineFragmentRandomItemTopic includes the requested fields of the GraphQL type Topic.
type SimpleInlineFragmentRandomItemTopic struct {
	Typename string `json:"__typename"`
	// ID is the identifier of the content.
	Id   testutil.ID `json:"id"`
	Name string      `json:"name"`
}

// SimpleInlineFragmentRandomItemVideo includes the requested fields of the GraphQL type Video.
type SimpleInlineFragmentRandomItemVideo struct {
	Typename string `json:"__typename"`
	// ID is the identifier of the content.
	Id       testutil.ID `json:"id"`
	Name     string      `json:"name"`
	Duration int         `json:"duration"`
}

// SimpleInlineFragmentResponse is returned by SimpleInlineFragment on success.
type SimpleInlineFragmentResponse struct {
	RandomItem SimpleInlineFragmentRandomItemContent `json:"-"`
}

func (v *SimpleInlineFragmentResponse) UnmarshalJSON(b []byte) error {

	type SimpleInlineFragmentResponseWrapper SimpleInlineFragmentResponse

	var firstPass struct {
		*SimpleInlineFragmentResponseWrapper
		RandomItem json.RawMessage `json:"randomItem"`
	}
	firstPass.SimpleInlineFragmentResponseWrapper = (*SimpleInlineFragmentResponseWrapper)(v)

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	{
		target := &v.RandomItem
		raw := firstPass.RandomItem
		err = __unmarshalSimpleInlineFragmentRandomItemContent(
			target, raw)
		if err != nil {
			return fmt.Errorf(
				"Unable to unmarshal SimpleInlineFragmentResponse.RandomItem: %w", err)
		}
	}
	return nil
}

func SimpleInlineFragment(
	client graphql.Client,
) (*SimpleInlineFragmentResponse, error) {
	var retval SimpleInlineFragmentResponse
	err := client.MakeRequest(
		nil,
		"SimpleInlineFragment",
		`
query SimpleInlineFragment {
	randomItem {
		__typename
		id
		name
		... on Article {
			text
		}
		... on Video {
			duration
		}
	}
}
`,
		&retval,
		nil,
	)
	return &retval, err
}

